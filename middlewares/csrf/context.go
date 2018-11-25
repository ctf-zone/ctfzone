package csrf

import (
	"context"
	"encoding/base64"
	"net/http"
)

type contextKey int

const (
	tokenKey contextKey = iota
)

func Token(r *http.Request) string {
	return r.Context().Value(tokenKey).(string)
}

func addTokenToRequest(r *http.Request, token []byte) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), tokenKey, base64.URLEncoding.EncodeToString(token)))
}
