package main

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/strata-io/service-extension/orchestrator"
)

// LoadAttrs derives a stable pseudonymous identifier for the authenticated user
// and tags the session with a fresh correlation ID, so an upstream application
// can recognise a returning user and support can trace a single visit without
// either value carrying the user's email.
//
// The two values look similar and are produced very differently. The pseudonym
// must be reproducible, so it comes from a hash of the email. The correlation
// ID must not be, so it comes from a random UUID. The FIPS-approved choice
// differs accordingly.
func LoadAttrs(api orchestrator.Orchestrator, _ http.ResponseWriter, _ *http.Request) error {
	logger := api.Logger()
	session, err := api.Session()
	if err != nil {
		return fmt.Errorf("unable to retrieve session: %w", err)
	}

	metadata := api.Metadata()
	// The namespace keeps a pseudonym scoped to one application, so the same
	// user's identifiers cannot be correlated across two of them.
	namespace := metadata["pseudonymNamespace"].(string)

	email, err := session.GetString("azure.email")
	if err != nil {
		return fmt.Errorf("failed to find user email: %w", err)
	}

	// The reflex here is crypto/md5 or crypto/sha1, and both are withdrawn from
	// the Service Extension runtime in FIPS 140-3 approved mode. They are also
	// the worst-behaved of the withdrawn packages: their Sum functions panic
	// rather than return an error when the validated cryptographic module is
	// enforcing, and a panicking extension disrupts the Orchestrator's request
	// processing. SHA-256 is approved and substitutes directly, needing only a
	// wider field for the longer digest.
	sum := sha256.Sum256([]byte(namespace + ":" + email))
	pseudonym := base64.RawURLEncoding.EncodeToString(sum[:])

	// uuid.NewMD5 and uuid.NewSHA1 mint name-based v3 and v5 UUIDs, and are
	// withdrawn along with the hashes they wrap. Neither is the right tool here
	// anyway: a correlation ID should not be derivable from the name it stands
	// for. uuid.New draws a random v4 UUID from crypto/rand, which is inside the
	// validated module. Use uuid.NewRandom if you would rather handle an entropy
	// failure as an error than take uuid.New's panic.
	correlationID := uuid.New().String()

	logger.Info(
		"se", "derived session identifiers",
		"correlationID", correlationID,
	)

	err = session.SetString("se.pseudonym", pseudonym)
	if err != nil {
		return fmt.Errorf("unable to set 'se.pseudonym' in session: %w", err)
	}
	err = session.SetString("se.correlationID", correlationID)
	if err != nil {
		return fmt.Errorf("unable to set 'se.correlationID' in session: %w", err)
	}

	return session.Save()
}
