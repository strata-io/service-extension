package cache

import (
	"context"
	"time"
)

type Cache interface {
	// GetBytes returns the []byte for a given key. If the key does not exist, an
	// error will be returned.
	GetBytes(ctx context.Context, key string) ([]byte, error)

	// SetBytes adds a key and the corresponding []byte value the backing store.
	// If options are passed, they will be configured for the key. Any existing value
	// for the key will be replaced.
	SetBytes(ctx context.Context, key string, value []byte, opts ...Option) error
}

// Options contains Options for a given piece of data.
type Options struct {
	// Represents the Time-To-Live (TTL) for a given piece of data. When this
	// time elapses, the data will be deleted from the underlying store.
	TTL time.Duration
}

// Option is an option to pass to the Cache setters. Option can be used to add
// additional capabilities or modify the default behavior for how the data is stored.
type Option func(*Options)

// WithTTL can be used to add a Time-To-Live (TTL) to a given piece of data.
func WithTTL(ttl time.Duration) Option {
	return func(do *Options) {
		do.TTL = ttl
	}
}

// Constraints are the constraints for a Cache.
type Constraints struct {
	// Name of the cache.
	Name string

	// RawNamespace indicates whether to use the namespace exactly as provided
	// without applying the /servext prefix. This allows Service Extensions to
	// access cache keys written by external systems.
	RawNamespace bool
}

// Constraint allows for customizing the Cache.
type Constraint func(*Constraints)

// WithName is an option to specify the cache name.
func WithName(name string) Constraint {
	return func(do *Constraints) {
		do.Name = name
	}
}

// WithRawNamespace enables raw namespace mode, bypassing the /servext prefix.
// Use this when accessing cache keys written by external systems that don't
// follow Maverics namespace conventions.
// FIXME: we may not want to do "`WithRawNamespace`" constraint here, instead we may just want an option which omits the `/servext` prefix.
func WithRawNamespace() Constraint {
	return func(do *Constraints) {
		do.RawNamespace = true
	}
}
