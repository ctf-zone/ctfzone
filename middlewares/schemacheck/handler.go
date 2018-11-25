package schemacheck

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/xeipuuv/gojsonschema"
)

type handler struct {
	next http.Handler
	opts *options
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Load request body.
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.opts.errorFunc(w, r, &Error{Msg: "invalid json"})
		return
	}

	// Validate schema.
	res, err := h.opts.schema.Validate(gojsonschema.NewBytesLoader(body))
	if err != nil {
		h.opts.errorFunc(w, r, &Error{Msg: "invalid json"})
		return
	}

	// Check validation successful.
	if !res.Valid() {
		e := &Error{Msg: "validation failed"}
		for _, err := range res.Errors() {
			e.Errs = append(e.Errs, &FieldError{Field: err.Field(), Msg: err.Description()})
		}

		h.opts.errorFunc(w, r, e)
		return
	}

	// Restore body
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	h.next.ServeHTTP(w, r)
}

func defaultErrorFunc(w http.ResponseWriter, r *http.Request, err error) {
	http.Error(w, "Bad request", 400)
}

func New(schema gojsonschema.JSONLoader, opts ...Option) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		do := *defaultOptions

		// Load schema
		s, err := gojsonschema.NewSchema(schema)

		if err != nil {
			panic(err)
		}

		do.schema = s

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
