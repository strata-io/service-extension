package orchestrator

import (
	"context"

	"github.com/strata-io/service-extension/app"
	"github.com/strata-io/service-extension/bundle"
	"github.com/strata-io/service-extension/cache"
	"github.com/strata-io/service-extension/http"
	"github.com/strata-io/service-extension/idfabric"
	"github.com/strata-io/service-extension/log"
	"github.com/strata-io/service-extension/router"
	"github.com/strata-io/service-extension/secret"
	"github.com/strata-io/service-extension/session"
	"github.com/strata-io/service-extension/tai"
	"github.com/strata-io/service-extension/weblogic"
)

type Orchestrator interface {
	// Logger gets a logger.
	Logger(opts ...log.Option) log.Logger

	// Session returns the session.
	Session(opts ...session.SessionOpt) (session.Provider, error)

	// SecretProvider gets a secret provider. An error is returned if a secret
	// provider is not configured.
	SecretProvider() (secret.Provider, error)

	// IdentityProvider gets an identity provider by name. An error is returned if
	// the identity provider is not found.
	IdentityProvider(name string) (idfabric.IdentityProvider, error)

	// AttributeProvider gets an attribute provider by name. An error is returned if
	// the attribute provider is not found.
	AttributeProvider(name string) (idfabric.AttributeProvider, error)

	// Metadata gets the metadata associated with the Service Extension in use.
	Metadata() map[string]any

	// Router gets a router.
	Router() router.Router

	// App gets the App associated with the Service Extension in use.
	App() (app.App, error)

	// TAI gets a TAI provider.
	TAI() tai.Provider

	// WebLogic gets a WebLogic provider.
	WebLogic() weblogic.Provider

	// Context gets the context associated with the Service Extension in use.
	// This is an experimental feature and may not be available in all Service Extensions.
	// If context is unavailable, nil will be returned.
	Context() context.Context

	// WithContext returns a shallow copy of an Orchestrator with the provided
	// context.
	WithContext(ctx context.Context) Orchestrator

	// Cache returns a cache that can be used to store state across different service
	// extensions.
	Cache(namespace string, opts ...cache.Constraint) (cache.Cache, error)

	// ServiceExtensionAssets exposes any assets that may have been bundled with the
	// service extension.
	ServiceExtensionAssets() bundle.SEAssets

	// HTTP provides utilities for making HTTP requests.
	HTTP() http.HTTP
}
