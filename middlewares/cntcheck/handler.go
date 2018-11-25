package cntcheck

import (
	"errors"
	"mime"
	"net/http"
	"strings"
)

type handler struct {
	next http.Handler
	opts *options
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mediatype, params, _ := mime.ParseMediaType(r.Header.Get("Content-Type"))
	charset, ok := params["charset"]
	if !ok {
		charset = "UTF-8"
	}

	// per net/http doc, means that the length is known and non-null
	if r.ContentLength > 0 &&
		!(mediatype == "application/json" && strings.ToUpper(charset) == "UTF-8") {
		h.opts.errorFunc(w, r, errors.New("bad Content-Type or charset"))
		return
	}

	h.next.ServeHTTP(w, r)
}

func defaultErrorFunc(w http.ResponseWriter, r *http.Request, err error) {
	http.Error(w, "Bad Content-Type or charset, expected 'application/json'", 415)
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
