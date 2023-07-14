package orchestrator

import (
	"github.com/strata-io/service-extension/log"
	"github.com/strata-io/service-extension/secret"
	"github.com/strata-io/service-extension/session"
)

type Orchestrator interface {
	Logger() log.Logger
	SessionProvider() session.Provider
	SecretProvider() (secret.Provider, error)
}
