package log

// Logger exposes three different levels of logging.
//
// The logger expects key-value pairs to enable structured logging.
//
// Example:
//
//	log.Info("msg", "example log", "requestID", "1234")
type Logger interface {
	// Debug will log at debug level.
	Debug(keyPairs ...any)
	// Info will log at info level.
	Info(keyPairs ...any)
	// Error will log at error level.
	Error(keyPairs ...any)
}
