# LDAP Search Service Extension (FIPS approved mode)

This is the approved-mode counterpart to the [`ldap-search`](../../ldap-search)
example. Both retrieve a user's group memberships from a directory and store
them on the session so a header rule can forward them upstream. They differ only
in how they reach LDAP.

`ldap-search` imports `github.com/go-ldap/ldap/v3` and does the work itself:
`ldap3.DialURL`, `conn.StartTLS`, `conn.Bind`, `conn.Search`. When the
Orchestrator runs in FIPS 140-3 approved mode those symbols are withdrawn from
the Service Extension runtime, because `*ldap.Conn` carries bind methods that
reach MD5 and MD4. An extension importing the package fails to load with an
error naming FIPS.

This example calls `api.AttributeProvider("ldapAttrs")` instead. The dial, the
TLS profile, and the service account bind become connector configuration in
[maverics.yaml](maverics.yaml), handled by the Orchestrator using approved
transport. The extension itself performs no cryptography, which is why it loads
and behaves identically in a standard and an approved deployment.

Note the other consequence of the move: the LDAP filter is no longer written in
Go. `Query` takes a subject and a list of attribute names, and the connector's
`usernameSearchKey` determines how the subject is matched. If you need a filter
the connector cannot express -- the multi-entry search the non-FIPS example was
written for -- see [Limits](#limits) below.

For the full list of what approved mode withdraws and what replaces it, see the
[FIPS examples README](../README.md).

## Setup

Please reference the [maverics.yaml](maverics.yaml) configuration file and the
Service Extension files it references for a set of action items specified with
`TODO`. These action items are changes necessary to get this example running.

The connector name in `metadata.attributeProviderName` must match the `name` of
the LDAP connector, and `metadata.groupAttribute` must name an attribute the
directory actually returns.

## Testing

1. Complete all the action items specified by `TODO`s.
1. Restart the Orchestrator, and ensure it starts successfully.
1. Navigate to the URL the Orchestrator is listening on in your browser:
   e.g. https://localhost/headers.
1. You should now be redirected to your specified IDP and prompted for
   authentication.
1. After successfully logging in, the `EXAMPLE-GROUPS` header will carry the
   groups retrieved through the attribute provider.

To confirm the FIPS behaviour, run the same test against an Orchestrator in
approved mode. This extension loads unchanged; the non-FIPS `ldap-search`
extension is refused at config load.

## Limits

`Query` covers attribute lookup for a known subject, which is what most
`loadAttrsSE` extensions need. It does not expose arbitrary LDAP searches --
multi-entry result sets, custom filters, or paged searches. If your integration
needs one of those in approved mode, model the lookup as an attribute on the
connector, or source the data over HTTP from a directory API. Reaching back to
`go-ldap` is not an option in approved mode.
