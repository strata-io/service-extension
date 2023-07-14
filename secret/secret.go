package secret

type Provider interface {
	// Get retrieves the key from the secret provider.
	Get(key string) any
	// GetString retrieves the key from the secret provider as a string value.
	GetString(key string) string
}
