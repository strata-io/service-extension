package cache

type Cache interface {
	// GetBytes returns the []byte for a given key. If the key does not exist, a
	// ErrNotFound will be returned.
	GetBytes(key string) ([]byte, error)

	// SetBytes adds a key and the corresponding []byte value the backing store.
	// If options are passed, they will be configured for the key. Any existing value
	// for the key will be replaced.
	SetBytes(key string, value []byte, opts ...Option) error
}

// Options contains Options for a given piece of data.
type Options struct {
	// Name of the cache.
	Name string
}

// Option is an option to pass to the Cache setters. Option can be used to add
// additional capabilities or modify the default behavior for how the data is stored.
type Option func(*Options)

// WithName is an option that allows the ability to specify the cache name.
func WithName(name string) Option {
	return func(do *Options) {
		do.Name = name
	}
}
