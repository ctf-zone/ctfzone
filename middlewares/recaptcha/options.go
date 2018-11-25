package recaptcha

import "net/http"

var defaultOptions = &options{
	header:    "X-G-Recaptcha-Response",
	secret:    "",
	errorFunc: defaultErrorFunc,
}

type options struct {
	header    string
	secret    string
	errorFunc func(http.ResponseWriter, *http.Request, error)
}

// Option defines the functional arguments for configuring the middleware.
type Option func(*options)

// Header sets the header name from which the recaptcha response will be taken.
// Default value is "X-G-Recaptcha-Response".
func Header(s string) Option {
	return func(opts *options) {
		opts.header = s
	}
}

// ErrorFunc allows you to control behavior when recaptcha check failed.
// The default behavior is for a HTTP 403 status code to be written to the
// ResponseWriter along with the plain-text error string.
func ErrorFunc(f func(http.ResponseWriter, *http.Request, error)) Option {
	return func(opts *options) {
		opts.errorFunc = f
	}
}
