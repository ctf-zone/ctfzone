package timecheck

import (
	"errors"
	"net/http"
	"time"
)

var (
	MinTime = time.Time{}
	MaxTime = time.Unix(1<<63-62135596801, 999999999)
)

var (
	ErrTooEarly = errors.New("timecheck: too early")
	ErrTooLate  = errors.New("timecheck: too late")
)

type handler struct {
	next http.Handler
	opts *options
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	now := time.Now().UTC()

	if now.Before(h.opts.start) {
		h.opts.errorFunc(w, r, ErrTooEarly)
		return
	}

	if now.After(h.opts.end) {
		h.opts.errorFunc(w, r, ErrTooLate)
		return
	}

	h.next.ServeHTTP(w, r)
}

func defaultErrorFunc(w http.ResponseWriter, r *http.Request, err error) {
	http.Error(w, "Access is limited for a period of time", 403)
}

func New(start, end time.Time, opts ...Option) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		do := *defaultOptions

		do.start = start
		do.end = end

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
