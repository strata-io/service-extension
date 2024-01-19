package orchestrator

import (
	"github.com/strata-io/service-extension/cache"
	"github.com/strata-io/service-extension/idfabric"
	"github.com/strata-io/service-extension/log"
	"github.com/strata-io/service-extension/router"
	"github.com/strata-io/service-extension/secret"
	"github.com/strata-io/service-extension/session"
	"github.com/strata-io/service-extension/tai"
)

type Orchestrator interface {
	// Logger gets a logger.
	Logger() log.Logger

	// Session returns the session.
	Session(opts ...session.SessionOpt) (session.Provider, error)

	// RequestCache returns a cache specifically to only be used by service
	// extensions.
	RequestCache(namespace string, opts ...cache.Constraint) (cache.Cache, error)

	// SecretProvider gets a secret provider. An error is returned if a secret
	// provider is not configured.
	SecretProvider() (secret.Provider, error)

	// IdentityProvider gets an identity provider by name. An error is returned if
	// the identity provider is not found.
	IdentityProvider(name string) (idfabric.IdentityProvider, error)

	// AttributeProvider gets an attribute provider by name. An error is returned if
	// the attribute provider is not found.
	AttributeProvider(name string) (idfabric.AttributeProvider, error)

	// Router gets a router.
	Router() router.Router

	// Metadata gets the metadata associated with the Service Extension in use.
	Metadata() map[string]any

	// TAI gets a TAI provider.
	TAI() tai.Provider
}
