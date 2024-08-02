package api

import "net/http"

type options struct {
	baseURL          string
	errorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// Option is a functional option for configuring the server.
type Option func(*options)

// WithBaseURL sets the base URL for the server.
func WithBaseURL(baseURL string) Option {
	return func(o *options) {
		o.baseURL = baseURL
	}
}

// WithErrorHandlerFunc sets the error handler function for the server.
func WithErrorHandlerFunc(f func(w http.ResponseWriter, r *http.Request, err error)) Option {
	return func(o *options) {
		o.errorHandlerFunc = f
	}
}
