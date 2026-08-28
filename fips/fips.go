// Package fips provides approved-mode-safe equivalents for the third-party
// libraries that are unavailable to service extensions when the Orchestrator
// runs in FIPS 140-3 approved mode.
//
// In approved mode the Orchestrator withdraws the go-jose and go-ldap symbols
// from the service extension runtime, because both can reach primitives the
// validated cryptographic module refuses. This package is the supported
// replacement for the JOSE half. The LDAP half is served by the existing
// attribute provider, reached through Orchestrator.AttributeProvider.
//
// The signature algorithm is derived from the key rather than chosen by the
// caller, so a non-approved algorithm cannot be named. The behaviour is
// identical in standard and FIPS builds, so an extension written against this
// package is portable across both and needs no build-time branching.
//
// Obtain a Provider from the Orchestrator:
//
//	signer, err := api.FIPS().NewSigner(privateKeyPEM)
//	if err != nil {
//		return fmt.Errorf("unable to create signer: %w", err)
//	}
//	token, err := signer.Sign(map[string]any{"sub": "alice"})
package fips

// Provider grants access to FIPS-approved cryptographic operations.
type Provider interface {
	// NewSigner returns a Signer over a PEM-encoded private key. The signature
	// algorithm follows the key type: RSA keys sign with RS256, and ECDSA keys
	// sign with ES256, ES384, or ES512 according to their curve. An error is
	// returned if the key cannot be parsed, or if its type or size is not
	// approved.
	NewSigner(privateKeyPEM []byte) (Signer, error)

	// NewVerifier returns a Verifier over either a JSON Web Key Set document or
	// a PEM-encoded public key. An error is returned if no usable key is found.
	NewVerifier(keys []byte) (Verifier, error)

	// Enabled reports whether the Orchestrator is running in FIPS 140-3
	// approved mode. An extension that must behave differently in an approved
	// deployment should branch on this rather than inspect the environment.
	Enabled() bool
}

// Signer mints signed JSON Web Tokens.
type Signer interface {
	// Sign returns the compact serialization of a JWT carrying claims. Any
	// registered claim the caller sets, such as "exp" or "iss", is preserved as
	// given; none are added automatically.
	Sign(claims map[string]any) (string, error)
}

// Verifier validates signed JSON Web Tokens.
type Verifier interface {
	// Verify checks the token's signature and returns its claims. It rejects a
	// token whose "alg" header falls outside the approved set, including the
	// "none" algorithm, before any signature check is attempted.
	//
	// Verify checks the signature only. Claim-level validation, such as
	// expiry, issuer, and audience, is the caller's responsibility.
	Verify(token string) (map[string]any, error)
}
