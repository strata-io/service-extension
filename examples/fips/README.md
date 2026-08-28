# Service Extensions in FIPS 140-3 Approved Mode

When the Maverics Identity Orchestrator runs in FIPS 140-3 approved mode, it
withdraws a set of symbols from the Service Extension runtime. Every one of them
can reach cryptography the validated cryptographic module refuses, and an
extension that imports a withdrawn package is refused at config load with an
error naming FIPS.

These examples show the approved way to do the three things authors most often
reach a withdrawn package for.

| Example | Replaces |
| --- | --- |
| [`ldap-search`](ldap-search) | `github.com/go-ldap/ldap/v3` |
| [`jwt-sign-verify`](jwt-sign-verify) | `github.com/go-jose/go-jose/v3` |
| [`hashing`](hashing) | `crypto/md5`, `crypto/sha1`, `uuid.NewMD5`, `uuid.NewSHA1` |

## What is withdrawn, and what to use instead

| Withdrawn | Approved alternative |
| --- | --- |
| `github.com/go-jose/go-jose/v3`, `.../v3/jwt` | `api.FIPS()` -- see the [`fips`](../../fips) package. The signature algorithm is derived from the key, so a non-approved algorithm cannot be named. |
| `github.com/go-ldap/ldap/v3` | `api.AttributeProvider("<connector name>")`. The dial, the TLS profile, and the service account bind become connector configuration, handled by the Orchestrator over approved transport. |
| `crypto/md5`, `crypto/sha1` | `crypto/sha256`, or any other SHA-2 or SHA-3 hash. |
| `crypto/des`, `crypto/rc4` | `crypto/aes`. |
| `crypto/dsa` | ECDSA via `crypto/ecdsa`, or RSA via `crypto/rsa`. |
| `uuid.NewMD5`, `uuid.NewSHA1` | `uuid.New` -- a random v4 UUID, drawn from `crypto/rand`. |
| `github.com/decred/dcrd/dcrec/secp256k1/v4` | No replacement. secp256k1 is a Koblitz curve outside the NIST set, so there is no approved equivalent. An integration that requires it cannot run in approved mode. |

Everything else stays available, including `crypto/rand`, `crypto/tls`,
`crypto/x509`, `crypto/sha256`, `crypto/ecdsa`, `crypto/rsa`, `crypto/ed25519`,
the `encoding/*` packages, and `net/http`.

### Why these five behave differently under an enforcing module

The withdrawals are not uniform in what they prevent, and the difference is
worth knowing when you are deciding how urgently to migrate an extension.

- `crypto/md5` and `crypto/sha1` **panic** inside `Sum` when the module is
  enforcing, and so do `uuid.NewMD5` and `uuid.NewSHA1`, which wrap them. A
  panicking extension disrupts the Orchestrator's request processing for that
  hook point. Withdrawing the symbol turns a mid-request crash into a
  config-load error.
- `crypto/des`, `crypto/rc4`, and `crypto/dsa` return a clean error. The
  withdrawal makes the failure happen at config load rather than at the first
  request that reaches the code path.
- `go-ldap` and `go-jose` are mixed. Some paths panic -- `go-ldap`'s `MD5Bind`
  and `NTLMBind`, `go-jose`'s JWKS parsing when a key carries an `x5c` chain --
  and some fail cleanly.
- `secp256k1` produces **no signal at all**. It implements its own constant-time
  field arithmetic and never enters the validated module, so an extension can
  generate keys and perform ECDH on a non-approved curve with nothing to
  indicate anything is wrong. Withdrawing the symbols is the only thing that
  closes it.

## Portability

An extension written against these APIs runs unchanged in both a standard and an
approved deployment. `api.FIPS()` and `api.AttributeProvider` are present and
behave identically in both, and none of the approved primitives above are
conditional. There is no build-time branching to do, and no separate copy of an
extension to maintain.

If an extension genuinely must behave differently in an approved deployment, use
`api.FIPS().Enabled()`. That is the supported way to detect approved mode --
prefer it over inspecting environment variables or build tags, which describe
how the binary was built rather than how it is running.

```go
if api.FIPS().Enabled() {
	logger.Info("se", "running in FIPS 140-3 approved mode")
}
```

## Where this restriction stops

Withdrawing symbols removes **packages** from what an extension can import. It
does not, and cannot, inspect what an extension computes.

An extension can still implement a non-approved primitive itself, in pure Go,
inside the interpreter -- MD5 from the RFC, a custom curve, a hand-rolled
cipher -- using nothing but arithmetic and `encoding/binary`. Nothing about that
code imports a withdrawn package, so nothing about this restriction sees it. The
same is true of an extension that ships a non-approved primitive as bundled
data, or reaches an external service that performs one on its behalf.

This is a real limit, honestly stated: the runtime enforces the **import
boundary**, not the algorithm. Extension-authored cryptography is covered by
review of the extension source before it is deployed, not by anything the
Orchestrator checks at load time. If you are responsible for an approved-mode
deployment, that review is a control you own.
