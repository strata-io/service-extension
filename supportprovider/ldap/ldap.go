package ldap

import (
	"crypto/tls"

	"github.com/go-ldap/ldap/v3"
)

// Scope choices.
const (
	ScopeBaseObject   = 0
	ScopeSingleLevel  = 1
	ScopeWholeSubtree = 2
)

// ScopeMap contains human-readable descriptions of scope choices.
var ScopeMap = map[int]string{
	ScopeBaseObject:   "Base Object",
	ScopeSingleLevel:  "Single Level",
	ScopeWholeSubtree: "Whole Subtree",
}

// Deref aliases.
const (
	NeverDerefAliases   = 0
	DerefInSearching    = 1
	DerefFindingBaseObj = 2
	DerefAlways         = 3
)

// DerefMap contains human-readable descriptions of deref aliases choices.
var DerefMap = map[int]string{
	NeverDerefAliases:   "NeverDerefAliases",
	DerefInSearching:    "DerefInSearching",
	DerefFindingBaseObj: "DerefFindingBaseObj",
	DerefAlways:         "DerefAlways",
}

// Support provides support of LDAP protocol.
type Support interface {
	DialURL(addr string, opts ...ldap.DialOpt) (*ldap.Conn, error)
	DialWithTLSConfig(tc *tls.Config) ldap.DialOpt
	NewSearchRequest(
		BaseDN string,
		Scope, DerefAliases, SizeLimit, TimeLimit int,
		TypesOnly bool,
		Filter string,
		Attributes []string,
		Controls []ldap.Control,
	) *ldap.SearchRequest
	NewPasswordModifyRequest(
		userIdentity string,
		oldPassword string,
		newPassword string,
	) *ldap.PasswordModifyRequest
}
