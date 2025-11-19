// Package ldap provides LDAP-specific interfaces for service extensions.
package ldap

import "context"

// BindDN retrieves the bind DN from the context. The context should be
// retrieved via the Orchestrator's Context() method (e.g., `api.Context()`).
// It returns an empty string if no bind DN is present.
type BindDN func(ctx context.Context) string
