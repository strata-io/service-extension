package cache

import "time"

type Cache interface {
	// GetBytes returns the []byte for a given key. If the key does not exist, an
	// error will be returned.
	GetBytes(key string) ([]byte, error)

	// SetBytes adds a key and the corresponding []byte value the backing store.
	// If options are passed, they will be configured for the key. Any existing value
	// for the key will be replaced.
	SetBytes(key string, value []byte, opts ...Option) error
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

// WithTTL can be used to add a Time-To-Live (TTL) to a given piece of data. When
// this option is set, it is up to the Store implementations to determine how
// best to honor this option.
func WithTTL(ttl time.Duration) Option {
	return func(do *Options) {
		do.TTL = ttl
	}
}

// Constraints contains Constraints for a given piece of data.
type Constraints struct {
	// Name of the cache.
	Name string
}

// Constraint is an option to pass to the Cache setters. Option can be used to add
// additional capabilities or modify the default behavior for how the data is stored.
type Constraint func(*Constraints)

// WithName is an option that allows the ability to specify the cache name.
func WithName(name string) Constraint {
	return func(do *Constraints) {
		do.Name = name
	}
}
