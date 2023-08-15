package secret

// Provider is used to retrieve secrets from the configured secret store.
type Provider interface {
	// Get retrieves the key from the secret provider.
	Get(key string) any

	// GetString retrieves the key from the secret provider as a string value.
	GetString(key string) string
}
