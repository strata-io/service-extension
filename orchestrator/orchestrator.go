package orchestrator

import (
	"github.com/strata-io/service-extension/log"
	"github.com/strata-io/service-extension/secret"
)

type Orchestrator interface {
	SecretProvider() (secret.Provider, error)
	Logger() log.Logger
}
