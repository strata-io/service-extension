// Package ldap provides LDAP-specific interfaces for service extensions.
package ldap

import "context"

// Provider is the interface for LDAP operations.
type Provider interface {
	// BindDN retrieves the bind DN from the context. The context should be
	// retrieved via the Orchestrator's Context() method (e.g., `api.Context()`).
	// It returns an empty string if no bind DN is present.
	BindDN(ctx context.Context) string
}
