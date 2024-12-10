package idfabric

import "net/http"

const (
	GrantTypeCliendCredentials = iota + 1
)

// IdentityProvider enables a way to interact with the identity provider.
// Interactions may include login and logout.
type IdentityProvider interface {
	// Login provides a front-channel user login flow. The user will be redirected
	// to the underlying IDP to authenticate the user.
	Login(rw http.ResponseWriter, req *http.Request, opts ...LoginOpt)
}

// LoginOptions store the options used to customize the user experience when
// calling Login on an IdentityProvider.
type LoginOptions struct {
	Username    string
	RedirectURL string

	GrantType int

	// ClientCredentialsResultCallback is called if GrantType is
	// 'GrantTypeClientCredentials' and it is not nil.
	// It is called at the end of Login() with the results of the
	ClientCredentialsResultCallback func(*TokenResult, *error)
}

// LoginOpt allows for customizing the login experience.
type LoginOpt func(cfg *LoginOptions)

// WithLoginHint specifies the username of the user to the IdentityProvider.
// This usually allows a known user to skip having to enter their username when
// prompted for authentication to the IdentityProvider.
func WithLoginHint(username string) LoginOpt {
	return func(cfg *LoginOptions) {
		cfg.Username = username
	}
}

// WithRedirectURL specifies landing page for the user after authenticating to
// the IdentityProvider.
func WithRedirectURL(url string) LoginOpt {
	return func(cfg *LoginOptions) {
		cfg.RedirectURL = url
	}
}

type TokenResult struct{}

// WithGrantTypeClientCredentials sets the grant type for this Login attempt to
// 'client_credentials'.
// The provided callback is called at the end of the Login() routine with the
// results.
func WithGrantTypeClientCredentials(callback func(t *TokenResult, e *error)) LoginOpt {
	return func(cfg *LoginOptions) {
		cfg.GrantType = GrantTypeCliendCredentials
		cfg.ClientCredentialsResultCallback = callback
	}
}

// AttributeProvider is used to retrieve attributes from an external system. A common
// attribute provider would be a data store such as LDAP.
type AttributeProvider interface {
	// Query is used to retrieve attributes for a user. A user's subject and the
	// requested attributes are consumed as params.
	//
	// When a query is successful, key-value pairs of the requested attributes are
	// returned. When a given AttributeProvider returns a multivalued attribute such
	// as group memberships, the values are concatenated using a delimiter that is
	// defined on the Identity Fabric component.
	Query(subject string, attributes []string) (map[string]string, error)
}
