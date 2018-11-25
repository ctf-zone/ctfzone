package csrf

import (
	"errors"
	"net/http"
)

var safeMethods = []string{"GET", "HEAD", "OPTIONS", "TRACE"}

type handler struct {
	next http.Handler
	opts *options
}

// contains returns true if string is in slice.
func contains(vals []string, s string) bool {
	for _, v := range vals {
		if v == s {
			return true
		}
	}
	return false
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	real := getRealToken(r, h.opts.cookie.name)
	if len(real) != tokenLen {
		token := genToken()
		sendNewToken(w, token, h.opts.cookie)
		r = addTokenToRequest(r, token)
	} else {
		r = addTokenToRequest(r, real)
	}

	if !contains(safeMethods, r.Method) {
		sent := getSentToken(r, h.opts.header)

		if real == nil || !compareTokens(real, sent) {
			h.opts.errorFunc(w, r, errors.New("tokens do not match"))
			return
		}
	}

	// Set the "Vary: Cookie" header to protect clients from caching the response.
	w.Header().Add("Vary", "Cookie")

	h.next.ServeHTTP(w, r)
}

func defaultErrorFunc(w http.ResponseWriter, r *http.Request, err error) {
	http.Error(w, "Wrong CSRF-token", 403)
}

func New(opts ...Option) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		do := *defaultOptions

		h := &handler{
			next: next,
			opts: &do,
		}

		for _, f := range opts {
			f(h.opts)
		}

		return h
	}
}
