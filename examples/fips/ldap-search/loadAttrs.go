package main

import (
	"fmt"
	"net/http"

	"github.com/strata-io/service-extension/orchestrator"
)

// LoadAttrs looks up a user's group memberships in LDAP and stores them on the
// session for later use.
//
// This is the approved-mode counterpart to the ldap-search example. That
// example dials the directory itself with github.com/go-ldap/ldap/v3, whose
// symbols the Orchestrator withdraws from the extension runtime in FIPS 140-3
// approved mode: the package's bind methods can reach MD5 and MD4, which the
// validated cryptographic module refuses.
//
// Reaching the directory through an attribute provider moves the connection,
// the TLS profile, and the service account bind into the Orchestrator's LDAP
// connector, which already enforces approved transport. Nothing in this file
// performs cryptography, so it behaves the same in a standard and an approved
// deployment.
func LoadAttrs(api orchestrator.Orchestrator, _ http.ResponseWriter, _ *http.Request) error {
	logger := api.Logger()
	session, err := api.Session()
	if err != nil {
		return fmt.Errorf("unable to retrieve session: %w", err)
	}

	metadata := api.Metadata()
	providerName := metadata["attributeProviderName"].(string)
	groupAttribute := metadata["groupAttribute"].(string)

	// The provider is resolved by connector name, so the directory URL, the
	// service account credentials, and the CA all stay in maverics.yaml rather
	// than being rebuilt in extension source on every request.
	attrProvider, err := api.AttributeProvider(providerName)
	if err != nil {
		return fmt.Errorf("unable to get attribute provider %q: %w", providerName, err)
	}

	uid, err := session.GetString("azure.email")
	if err != nil {
		return fmt.Errorf("failed to find user email required for LDAP query: %w", err)
	}

	logger.Info("se", "loading attributes from LDAP")

	attrs, err := attrProvider.Query(uid, []string{groupAttribute})
	if err != nil {
		return fmt.Errorf("unable to query attribute provider: %w", err)
	}

	// Query flattens a multivalued attribute into one string using the
	// 'attributeDelimiter' configured on the connector, so the groups arrive
	// already joined and need no delimiter of their own here.
	groups, ok := attrs[groupAttribute]
	if !ok {
		logger.Info(
			"se", "no groups returned for user",
			"attribute", groupAttribute,
		)
		return nil
	}

	logger.Debug(
		"se", "setting groups attribute on session",
		"se.groups", groups,
	)

	err = session.SetString("se.groups", groups)
	if err != nil {
		return fmt.Errorf("unable to set 'se.groups' in session: %w", err)
	}
	err = session.Save()
	if err != nil {
		return fmt.Errorf("unable to save session: %w", err)
	}

	return nil
}
