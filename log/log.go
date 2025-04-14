package log

import "net/http"

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

// Option is an option used to configure the retrieval of the Logger.
//
// Example:
//
//	logger := api.Logger(WithFromRequest(req))
type Option func(*Options)

// WithFromRequest retrieves the logger from the request's context. If no logger is
// found on the request's context or the request is nil, a default logger will be
// returned.
//
// This method only needs be used in service extensions that expose their own
// HTTP handlers using router.HandleFunc. Using this request bound logger ensures
// request specific key-value pairs (e.g. traceID) are included in log messages.
func WithFromRequest(r *http.Request) Option {
	return func(o *Options) {
		o.Request = r
	}
}

type Options struct {
	Request *http.Request
}
