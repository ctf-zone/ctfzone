package schemacheck

import (
	"net/http"

	"github.com/xeipuuv/gojsonschema"
)

var defaultOptions = &options{
	schema:    nil,
	errorFunc: defaultErrorFunc,
}

type options struct {
	schema    *gojsonschema.Schema
	errorFunc func(http.ResponseWriter, *http.Request, error)
}

// Option defines the functional arguments for configuring the middleware.
type Option func(*options)

// ErrorFunc allows you to control behavior when recaptcha check failed.
// The default behavior is for a HTTP 403 status code to be written to the
// ResponseWriter along with the plain-text error string.
func ErrorFunc(f func(http.ResponseWriter, *http.Request, error)) Option {
	return func(opts *options) {
		opts.errorFunc = f
	}
}
