package session

import (
	"net/http"
	"time"
)

// Provider enables a way to interact with the underlying session store. Methods
// on the provider take a request in order to look up the associated session.
//
// Example (basic setters and getters):
//
//	sess, _ := api.Session()
//	_ = sess.SetString("idp.authenticated", "true")
//	_ = sess.Save()
//	isAuthn, _ := sess.GetString("idp.authenticated")
//
// Example (JSON setters and getters):
//
//	sess, _ := api.Session()
//	type address struct {
//		City    string `json:"city"`
//		State   string `json:"state"`
//		Country string `json:"country"`
//	}
//
//	type person struct {
//		Name    string  `json:"name"`
//		Age     int     `json:"age"`
//		Address address `json:"address"`
//	}
//
//	var p = person{
//		Name: "John Doe",
//		Age:  42,
//		Address: address{
//			City:    "Boulder",
//			State:   "CO",
//			Country: "USA",
//		},
//	}
//	_ := sess.SetJSON("profile", &p)
//	_ = sess.Save()
//
//	var p2 person
//	_ := sess.GetJSON("profile", &p2)
type Provider interface {
	// GetString returns a session value based on the provided key. If the key does
	// not exist, the default or zero value will be returned (i.e, "").
	GetString(key string) (string, error)

	// GetBool returns a session value based on the provided key. If the key does
	// not exist, the default or zero value will be returned (i.e, false).
	GetBool(key string) (bool, error)

	// GetInt returns a session value based on the provided key. If the key does
	// not exist, the default or zero value will be returned (i.e, 0).
	GetInt(key string) (int, error)

	// GetFloat returns a session value based on the provided key. If the key does
	// not exist, the default or zero value will be returned (i.e, 0.0).
	GetFloat(key string) (float64, error)

	// GetBytes returns the []byte for a given key from the session data. If the key
	// does not exist, the default or zero value will be returned (i.e, nil).
	GetBytes(key string) ([]byte, error)

	// GetTime returns the time.Time for a given key from the session data. If the key
	// does not exist, the default or zero value will be returned
	// (i.e, 0001-01-01 00:00:00 +0000 UTC).
	GetTime(key string) (time.Time, error)

	// GetJSON parses the JSON-encoded data for a given key from the session and
	// stores the result in the value pointed to by dest. If dest is nil or not a
	// pointer, an error is returned.
	GetJSON(key string, dest any) error

	// SetString adds a key and the corresponding string value to the session data.
	SetString(key string, value string) error

	// SetInt adds a key and the corresponding int value to the session data.
	SetInt(key string, value int) error

	// SetFloat adds a key and the corresponding float value to the session data.
	SetFloat(key string, value float64) error

	// SetBool adds a key and the corresponding boolean value to the session data.
	SetBool(key string, value bool) error

	// SetBytes adds a key and the corresponding []byte value to the session data.
	SetBytes(key string, value []byte) error

	// SetTime adds a key and the corresponding time.Time value to the session data.
	SetTime(key string, value time.Time) error

	// SetJSON adds a key and the corresponding value to the session data. The
	// value must be JSON encodable or an error will be returned.
	SetJSON(key string, value any) error

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
