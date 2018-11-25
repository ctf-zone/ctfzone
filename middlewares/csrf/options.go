package csrf

import (
	"net/http"
	"time"
)

var defaultOptions = &options{
	header: "X-CSRF-Token",
	cookie: &cookieOptions{
		name:     "csrf-token",
		secure:   true,
		lifetime: time.Hour * 24,
		httponly: false,
	},
	errorFunc: defaultErrorFunc,
}

type cookieOptions struct {
	name     string
	secure   bool
	lifetime time.Duration
	httponly bool
}

type options struct {
	header    string
	cookie    *cookieOptions
	errorFunc func(http.ResponseWriter, *http.Request, error)
}

// Option defines the functional arguments for configuring the middleware.
type Option func(*options)

// Header sets the header name from which the CSRF-token will be taken.
// Default value is "X-CSRF-Token".
func Header(s string) Option {
	return func(opts *options) {
		opts.header = s
	}
}

func CookieName(s string) Option {
	return func(opts *options) {
		opts.cookie.name = s
	}
}

func CookieSecure(b bool) Option {
	return func(opts *options) {
		opts.cookie.secure = b
	}
}

func CookieLifetime(t time.Duration) Option {
	return func(opts *options) {
		opts.cookie.lifetime = t
	}
}

func CookieHTTPOnly(b bool) Option {
	return func(opts *options) {
		opts.cookie.httponly = b
	}
}

func ErrorFunc(f func(http.ResponseWriter, *http.Request, error)) Option {
	return func(opts *options) {
		opts.errorFunc = f
	}
}
