package supportprovider

import (
	"github.com/strata-io/service-extension/supportprovider/ldap"
)

// SupportProvider instance provides support to selected list of external
// services and protocols.
type SupportProvider interface {
	ldap.Support
}
