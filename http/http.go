package http

import (
	"net/http"

	"github.com/strata-io/service-extension/log"
)

// HTTP provides a way to interact with the HTTP client.
type HTTP interface {
	// GetClient returns the HTTP client based on the provided name.
	// If the client does not exist, an error will be returned.
	GetClient(name string) (*http.Client, error)

	// SetClient adds a client to the HTTP client store based on the provided name.
	SetClient(name string, client *http.Client) error

	// DefaultClient returns the default HTTP client.
	DefaultClient() *http.Client

	// LoggerFromRequest returns the logger from the request's context. If no logger
	// is found on the request's context or the request is nil, a default logger
	// will be returned.
	//
	// This method only needs be used in service extensions that expose their own
	// HTTP handlers using router.HandleFunc. Using this request bound logger ensures
	// request specific key-value pairs (e.g. traceID) are included in log messages.
	LoggerFromRequest(req *http.Request) log.Logger
}
