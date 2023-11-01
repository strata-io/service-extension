package session

import "net/http"

// Provider enables a way to interact with the underlying session store. Methods
// on the provider take a request in order to look up the associated session.
//
// Example:
//
//	sess, _ := api.Session()
//	_ = sess.SetString("idp.authenticated", "true")
//	_ = sess.Save()
type Provider interface {
	// GetString returns a session value based on the provided key.
	GetString(key string) (string, error)

	// SetString adds a key and the corresponding string value to the session data.
	SetString(key string, value string) error

	// Save saves all changes from the changelog to the underlying session store.
	Save() error
}

type Options struct {
	Request *http.Request
}

// SessionOpt is an option that allows to configure retrieval of the session.
//
// Example:
//
//	sess, _ := api.Session(WithRequest(req))
//	isAuth, _ := sess.GetString("idp.authenticated")
type SessionOpt func(*Options)

// WithRequest is a SessionOpt that allows to retrieve a particular session with
// a specific request.
//
// Example:
//
//	sess, _ := api.Session(WithRequest(req))
//	isAuth, _ := sess.GetString("idp.authenticated")
func WithRequest(req *http.Request) SessionOpt {
	return func(o *Options) {
		o.Request = req
	}
}
