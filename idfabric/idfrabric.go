package idfabric

import "net/http"

// IdentityProvider enables a way to interact with the identity provider.
// Interactions may include login and logout.
type IdentityProvider interface {
	// Login provides a front-channel user login flow. The user will be redirected
	// to the underlying IDP to authenticate the user.
	Login(rw http.ResponseWriter, req *http.Request)
}
