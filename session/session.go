package session

// Session represents the state of an end-user.
type Session interface {
	// GetString returns a session value based on the provided key.
	GetString(key string) (string, error)

	// SetString sets a value on the session for the provided key.
	SetString(key string, value any) error

	// Save saves all changes from the changelog to the underlying session store.
	Save() error
}
