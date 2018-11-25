package recaptcha

import (
	"net/http"
)

type handler struct {
	next http.Handler
	opts *options
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := check(
		h.opts.secret,
		r.Header.Get(h.opts.header),
		r.Header.Get("X-Forwared-For"))

	if err != nil {
		h.opts.errorFunc(w, r, err)
		return
	}

	h.next.ServeHTTP(w, r)
}

func defaultErrorFunc(w http.ResponseWriter, r *http.Request, err error) {
	http.Error(w, "Wrong captcha", 403)
}

func New(secret string, opts ...Option) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		do := *defaultOptions

		do.secret = secret

		h := &handler{
			next: next,
			opts: &do,
		}

		for _, fn := range opts {
			fn(h.opts)
		}

		return h
	}
}
