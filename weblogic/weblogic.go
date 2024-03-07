package weblogic

import (
	"time"
)

type Config struct {
	// RSAPrivateKeyPEM is the pem-encoded RSA PKCS1 private key that will be used to
	// sign the JWT.
	RSAPrivateKeyPEM string

	// Subject is the user's unique identifier. This value will be mapped to the
	// JWT's 'sub' claim.
	Subject string

	// Lifetime is the duration of the token's lifetime. This value will be mapped
	// to the JWT's 'exp' claim. The Lifetime should generally be set to match the
	// lifetime of a user's session.
	Lifetime time.Duration
}

type Provider interface {
	// NewSignedJWT returns a signed JWT that the WebLogic Identity Asserter module
	// will consume in order to build its identity context.
	NewSignedJWT(Config) (string, error)
}
