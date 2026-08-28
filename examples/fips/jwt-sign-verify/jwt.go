package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/strata-io/service-extension/log"
	"github.com/strata-io/service-extension/orchestrator"
)

// MintToken signs a short-lived JWT asserting the authenticated user's identity
// and stores it on the session, where a header rule can pick it up and forward
// it to the upstream application.
//
// In FIPS 140-3 approved mode the Orchestrator withdraws the
// github.com/go-jose/go-jose/v3 symbols from the Service Extension runtime,
// because go-jose names its algorithms with bare strings and so can be pointed
// at primitives the validated cryptographic module refuses. api.FIPS() is the
// supported replacement. It behaves the same in a standard and an approved
// deployment, so this extension needs no build-time branching.
func MintToken(api orchestrator.Orchestrator, _ http.ResponseWriter, _ *http.Request) error {
	logger := api.Logger()
	session, err := api.Session()
	if err != nil {
		return fmt.Errorf("unable to retrieve session: %w", err)
	}

	metadata := api.Metadata()
	issuer := metadata["issuer"].(string)
	audience := metadata["audience"].(string)

	// Read the signing key from the secret provider rather than embedding it in
	// the extension. Service Extension source travels in the config bundle as
	// plain text, so a key written here would be readable by anyone who can read
	// the bundle.
	secretProvider, err := api.SecretProvider()
	if err != nil {
		return fmt.Errorf("unable to get secret provider: %w", err)
	}
	privateKeyPEM := secretProvider.GetString("jwtSigningKey")
	if privateKeyPEM == "" {
		return errors.New("secret 'jwtSigningKey' is empty or not configured")
	}

	// The algorithm follows the key: an RSA key signs with RS256, and an ECDSA
	// key with ES256, ES384, or ES512 to match its curve. There is no algorithm
	// parameter, so a non-approved choice is not merely rejected at runtime --
	// it cannot be named.
	signer, err := api.FIPS().NewSigner([]byte(privateKeyPEM))
	if err != nil {
		return fmt.Errorf("unable to create signer: %w", err)
	}

	subject, err := session.GetString("azure.email")
	if err != nil {
		return fmt.Errorf("failed to find user email required for the subject claim: %w", err)
	}

	// Sign preserves the registered claims exactly as given and adds none of its
	// own, so a token has an expiry only because one is set here.
	now := time.Now()
	token, err := signer.Sign(map[string]any{
		"iss": issuer,
		"aud": audience,
		"sub": subject,
		"iat": now.Unix(),
		"exp": now.Add(5 * time.Minute).Unix(),
	})
	if err != nil {
		return fmt.Errorf("unable to sign token: %w", err)
	}

	logger.Info("se", "minted upstream token", "subject", subject)

	err = session.SetString("se.upstreamToken", token)
	if err != nil {
		return fmt.Errorf("unable to set 'se.upstreamToken' in session: %w", err)
	}

	return session.Save()
}

// ValidateToken verifies the bearer token on an inbound request and authorizes
// the request only if both the signature and the claims hold up.
func ValidateToken(api orchestrator.Orchestrator, _ http.ResponseWriter, req *http.Request) bool {
	logger := api.Logger(log.WithRequest(req))

	metadata := api.Metadata()
	issuer := metadata["issuer"].(string)

	secretProvider, err := api.SecretProvider()
	if err != nil {
		logger.Error("se", "unable to get secret provider", "error", err)
		return false
	}

	// NewVerifier accepts either a JSON Web Key Set document or a PEM-encoded
	// public key, so the same call works whether keys are fetched from an
	// issuer's JWKS endpoint or configured directly.
	verifier, err := api.FIPS().NewVerifier([]byte(secretProvider.GetString("jwtVerificationKeys")))
	if err != nil {
		logger.Error("se", "unable to create verifier", "error", err)
		return false
	}

	token, ok := bearerToken(req)
	if !ok {
		logger.Info("se", "request denied", "reason", "no bearer token on request")
		return false
	}

	claims, err := verifier.Verify(token)
	if err != nil {
		// A rejected token is ordinary traffic, not an operator problem: an
		// unknown signing key, a tampered token, or an 'alg' header outside the
		// approved set are all reported here.
		logger.Info("se", "request denied", "reason", "signature verification failed", "error", err)
		return false
	}

	// Verify checks the signature only. Every claim-level check -- expiry,
	// issuer, audience, scope -- is the caller's responsibility, so a token that
	// verifies is not yet a token that should be honoured.
	if !validClaims(api, claims, issuer) {
		return false
	}

	logger.Info("se", "request authorized", "subject", claims["sub"])
	return true
}

// validClaims applies the claim checks Verify deliberately leaves to the
// caller.
func validClaims(api orchestrator.Orchestrator, claims map[string]any, issuer string) bool {
	logger := api.Logger()

	if got, _ := claims["iss"].(string); got != issuer {
		logger.Info("se", "request denied", "reason", "unexpected issuer", "iss", got)
		return false
	}

	// Claims are decoded from JSON, so numeric claims such as 'exp' and 'iat'
	// arrive as float64 rather than as an integer type.
	exp, ok := claims["exp"].(float64)
	if !ok {
		logger.Info("se", "request denied", "reason", "token has no 'exp' claim")
		return false
	}
	if time.Now().After(time.Unix(int64(exp), 0)) {
		logger.Info("se", "request denied", "reason", "token is expired")
		return false
	}

	return true
}

// bearerToken extracts the credential from an Authorization header.
func bearerToken(req *http.Request) (string, bool) {
	const prefix = "Bearer "

	header := req.Header.Get("Authorization")
	if len(header) <= len(prefix) || !strings.EqualFold(header[:len(prefix)], prefix) {
		return "", false
	}

	return strings.TrimSpace(header[len(prefix):]), true
}
