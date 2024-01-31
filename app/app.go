package app

// App enables a way to interact with the application which defines the
// Service Extension.
type App interface {
	// Name returns the name of the application.
	Name() string
}
