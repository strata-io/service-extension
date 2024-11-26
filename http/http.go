package http

import "net/http"

// HTTP provides a way to interact with the HTTP client.
type HTTP interface {
	// GetClient returns the HTTP client based on the provided name.
	// If the client does not exist, the second return value will be false.
	GetClient(name string) (*http.Client, bool)

	// SetClient adds a client to the HTTP client store based on the provided name.
	SetClient(name string, client *http.Client)

	// DefaultClient returns the default HTTP client.
	DefaultClient() *http.Client
}
