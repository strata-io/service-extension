package idfabric

import (
	"net/http"
	"net/url"
)

const (
	GrantTypeROPC = iota + 1
)

// IdentityProvider enables a way to interact with the identity provider.
// Interactions may include login and logout.
type IdentityProvider interface {
	// Login provides a front-channel user login flow. The user will be redirected
	// to the underlying IDP to authenticate the user.
	Login(rw http.ResponseWriter, req *http.Request, opts ...LoginOpt)
	// IsAvailable checks if the underlying IDP is available for use.
	IsAvailable() bool
}

// LoginOptions store the options used to customize the user experience when
// calling Login on an IdentityProvider.
type LoginOptions struct {
	Username             string
	RedirectURL          string
	SilentAuthentication bool
	QueryParams          url.Values
	GrantType            int
	ROPCRequest          ROPCRequest
	LoginResult          *LoginResult
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

// WithQueryParam enables a way to specify custom query parameters to be added to the
// authorization request.
func WithQueryParam(k, v string) LoginOpt {
	return func(cfg *LoginOptions) {
		if len(cfg.QueryParams) == 0 {
			cfg.QueryParams = url.Values{}
		}
		cfg.QueryParams.Add(k, v)
	}
}

// WithSilentAuthentication specifies to the IDP that no user interaction should
// occur as part of the login.
//
// In the context of OIDC, this option will result in the 'prompt=none' query
// parameter being sent as part of the authentication request. For more details,
// please see the OIDC RFC
// https://openid.net/specs/openid-connect-core-1_0-final.html#AuthRequest.
func WithSilentAuthentication() LoginOpt {
	return func(cfg *LoginOptions) {
		cfg.SilentAuthentication = true
	}
}

// ROPCRequest is used to authenticate to the IdentityProvider using the
// Resource Owner Password Credentials (ROPC) flow.
type ROPCRequest struct {
	Username string
	Password string
}

// LoginResult is the response from the IdentityProvider after a login attempt.
type LoginResult struct {
	TokenResult
	Error error
}

// TokenResult is the response from the IdentityProvider after a login attempt.
type TokenResult struct {
	// AccessToken is the token that can be used to access protected resources.
	AccessToken string
	// IDToken is the token that contains user identity information.
	IDToken string
	// RefreshToken is the token that can be used to refresh the access token.
	RefreshToken string
	// ExpiresIn is the number of seconds until the access token expires.
	ExpiresIn int64
	// Scope is the scope of the access token.
	Scope string
}

// WithGrantTypeROPC specifies the Resource Owner Password Credentials (ROPC)
// flow for authenticating a user. This flow is typically used for legacy
// applications that require a username and password to authenticate the user
// directly. https://datatracker.ietf.org/doc/html/rfc6749#section-4.3
func WithGrantTypeROPC(input ROPCRequest, output *LoginResult) LoginOpt {
	return func(cfg *LoginOptions) {
		cfg.GrantType = GrantTypeROPC
		cfg.ROPCRequest = input
		cfg.LoginResult = output
	}
}

// WithRedirectURL specifies landing page for the user after authenticating to
// the IdentityProvider.
func WithRedirectURL(url string) LoginOpt {
	return func(cfg *LoginOptions) {
		cfg.RedirectURL = url
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
