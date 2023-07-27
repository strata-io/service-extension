package session

import "net/http"

// Provider enables a way to interact with the underlying session store. Methods
// on the provider take a request in order to lookup the associated session.
//
// Example:
//
//	session.Set(req, "idp.authenticated", "true")
type Provider interface {
	// GetString returns a session value based on the provided key.
	GetString(req *http.Request, key string) string

	// Get returns a session value based on the provided key.
	Get(req *http.Request, key string) any

	// Set sets a value on the session for the provided key.
	Set(req *http.Request, key string, value any)
}

// Session provides a session instance for an end-user. This interface is
// typically only used for back-channel interactions that do not contain a
// session cookie on the request.
type Session interface {
	// GetString returns a session value based on the provided key.
	GetString(key string) string

	// Get returns a session value based on the provided key.
	Get(key string) any

	// Set sets a value on the session for the provided key.
	Set(key string, value any)
}
