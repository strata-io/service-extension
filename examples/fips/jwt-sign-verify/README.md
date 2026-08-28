# JWT Sign and Verify Service Extension (FIPS approved mode)

This example mints a signed JWT for the authenticated user and forwards it to
the upstream application, then verifies an inbound bearer token on a second
route. Both halves go through `api.FIPS()`, the supported JOSE path for Service
Extensions.

In FIPS 140-3 approved mode the Orchestrator withdraws
`github.com/go-jose/go-jose/v3` and `.../v3/jwt` from the Service Extension
runtime. go-jose names its algorithms with bare strings, so any algorithm can be
requested, including ones the validated cryptographic module refuses; and
`JSONWebKey.UnmarshalJSON` reaches SHA-1 when a key carries an `x5c` chain,
which panics under an enforcing module rather than returning an error.

`api.FIPS()` closes both. The algorithm is derived from the key rather than
chosen by the caller, so a non-approved algorithm is inexpressible rather than
merely rejected, and JWKS parsing handles certificate chains without reaching a
non-approved hash. The behaviour is identical in a standard and an approved
build, so an extension written against it is portable and needs no build-time
branching.

For the full list of what approved mode withdraws and what replaces it, see the
[FIPS examples README](../README.md).

## What the example shows

**`MintToken`** signs a token:

- The signing key comes from `api.SecretProvider()`, not from the source.
  Service Extension source travels in the config bundle as plain text.
- `api.FIPS().NewSigner(pem)` picks the algorithm from the key. An RSA key signs
  with RS256; an ECDSA key with ES256, ES384, or ES512 to match its curve.
- `Sign` preserves the registered claims exactly as given and adds none, so the
  extension sets `iss`, `aud`, `sub`, `iat`, and `exp` itself.

**`ValidateToken`** verifies a token:

- `api.FIPS().NewVerifier(keys)` accepts either a JWKS document or a PEM-encoded
  public key.
- The error from `Verify` is handled rather than discarded. An unknown key, a
  tampered signature, and an `alg` header outside the approved set (including
  `none`) all surface here.
- **`Verify` checks the signature only.** Claim validation is the caller's job,
  which is why the example checks `iss` and `exp` explicitly after a successful
  verification. A token that verifies is not yet a token that should be
  honoured.

## Setup

Please reference the [maverics.yaml](maverics.yaml) configuration file and the
Service Extension files it references for a set of action items specified with
`TODO`, along with the keys in [secrets.yaml](secrets.yaml).

## Testing

1. Complete all the action items specified by `TODO`s.
1. Restart the Orchestrator, and ensure it starts successfully.
1. Navigate to the URL the Orchestrator is listening on in your browser and
   authenticate. The upstream request will carry an `Authorization: Bearer`
   header holding the minted token.
1. Send a request to `/api` with that token in an `Authorization` header. It
   should be allowed, and the Orchestrator should log
   `request authorized`.
1. Send the same request with a token whose signature or `exp` has been altered.
   It should be denied, and the Orchestrator should log which check failed.

Running the same test against an Orchestrator in approved mode should produce
identical results. That is the point of the package.
