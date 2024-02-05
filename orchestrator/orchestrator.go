package orchestrator

import (
	"context"
	"io/fs"

	"github.com/strata-io/service-extension/app"
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

	// Cache returns a cache that can be used to store state across different service
	// extensions.
	Cache(namespace string, opts ...cache.Constraint) (cache.Cache, error)

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

	// App gets the App associated with the Service Extension in use.
	App() (app.App, error)

	// Context gets the context associated with the Service Extension in use.
	// This is an experimental feature and may not be available in all Service Extensions.
	// If context is unavailable, nil will be returned.
	Context() context.Context

	// ConfigFS gets the configured ConfigReader.
	ConfigFS() ConfigReader
}

type ConfigReader interface { // TODO ORC: should this be in another package?
	FS() (fs.FS, error)
	//ReadFile(string) ([]byte, error)
}
