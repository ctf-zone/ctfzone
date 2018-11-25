package cntcheck

import (
	"net/http"
)

var defaultOptions = &options{
	errorFunc: defaultErrorFunc,
}

type options struct {
	errorFunc func(http.ResponseWriter, *http.Request, error)
}

// Option defines the functional arguments for configuring the middleware.
type Option func(*options)

func ErrorFunc(f func(http.ResponseWriter, *http.Request, error)) Option {
	return func(opts *options) {
		opts.errorFunc = f
	}
}
