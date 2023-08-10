package router

import "net/http"

// Router is used to register HTTP endpoints on the Orchestrator.
type Router interface {
	// HandleFunc registers the handler function for the given pattern in the
	// Router. An error is returned if a given route is already registered.
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) error
}
