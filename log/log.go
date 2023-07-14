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
	Debug(keyPairs ...string)
	// Info will log at info level.
	Info(keyPairs ...string)
	// Error will log at error level.
	Error(keyPairs ...string)
}
