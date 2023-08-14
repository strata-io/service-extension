package idfabric

import "net/http"

// IdentityProvider enables a way to interact with the identity provider.
// Interactions may include login and logout.
type IdentityProvider interface {
	// Login provides a front-channel user login flow. The user will be redirected
	// to the underlying IDP to authenticate the user.
	Login(rw http.ResponseWriter, req *http.Request)
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
