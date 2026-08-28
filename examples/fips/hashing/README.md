# Hashing and ID Generation Service Extension (FIPS approved mode)

This example derives two identifiers for the authenticated user and forwards
them upstream as headers: a stable pseudonym standing in for the user's email,
and a fresh correlation ID for the visit. Neither carries the email itself.

The two are the most common reasons a Service Extension reaches for a hash, and
they need opposite properties. The pseudonym must be reproducible, so the same
user maps to the same value on every request -- that requires a hash of the
email. The correlation ID must not be reproducible, and must not be derivable
from any identity it is attached to -- that requires randomness.

## What approved mode changes

| Reflex | Approved alternative |
| --- | --- |
| `md5.Sum` / `sha1.Sum` for a stable identifier | `sha256.Sum256`, or any SHA-2 / SHA-3 hash |
| `uuid.NewMD5` / `uuid.NewSHA1` for a name-based UUID | `uuid.New` for a random v4 UUID |

`crypto/md5` and `crypto/sha1` are the worst-behaved of the packages approved
mode withdraws. Their `Sum` functions **panic** rather than return an error when
the validated cryptographic module is enforcing, and a panicking extension
disrupts the Orchestrator's request processing. Withdrawing the symbols converts
that mid-request crash into a config-load error naming FIPS. `uuid.NewMD5` and
`uuid.NewSHA1` are withdrawn for the same reason -- they are thin wrappers over
those two hashes.

SHA-256 substitutes for MD5 or SHA-1 directly in the pseudonym case; the only
adjustment is a wider field for the longer digest. `uuid.New` is not a
substitute for `uuid.NewMD5` in general -- a v4 UUID is random where a v3 UUID
is a function of its name -- but for a correlation ID randomness is the property
you actually wanted. If you genuinely need a name-derived identifier, hash the
name with SHA-256 as this example does for the pseudonym, rather than looking
for a name-based UUID constructor that survives approved mode. There isn't one.

For the full list of what approved mode withdraws and what replaces it, see the
[FIPS examples README](../README.md).

## Setup

Please reference the [maverics.yaml](maverics.yaml) configuration file and the
Service Extension files it references for a set of action items specified with
`TODO`.

## Testing

1. Complete all the action items specified by `TODO`s.
1. Restart the Orchestrator, and ensure it starts successfully.
1. Navigate to the URL the Orchestrator is listening on in your browser:
   e.g. https://localhost/headers.
1. After successfully logging in, confirm the `EXAMPLE-USER-ID` and
   `EXAMPLE-CORRELATION-ID` headers are rendered on the `/headers` page of the
   sample app.
1. Log out and back in. `EXAMPLE-USER-ID` should be unchanged and
   `EXAMPLE-CORRELATION-ID` should be different.

This extension loads and behaves identically in a standard and an approved
deployment.
