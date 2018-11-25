package csrf

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"io"
	"net/http"
	"time"
)

var tokenLen = 32

func genToken() []byte {
	token := make([]byte, tokenLen)

	if _, err := io.ReadFull(rand.Reader, token); err != nil {
		panic(err)
	}

	return token
}

func compareTokens(a, b []byte) bool {
	return subtle.ConstantTimeCompare(a, b) == 1
}

func getSentToken(r *http.Request, header string) []byte {

	encoded := r.Header.Get(header)

	token, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		return nil
	}

	return token
}

func getRealToken(r *http.Request, cookieName string) []byte {

	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return nil
	}

	token, err := base64.URLEncoding.DecodeString(cookie.Value)
	if err != nil {
		return nil
	}

	return token
}

func sendNewToken(w http.ResponseWriter, token []byte, opts *cookieOptions) {

	cookie := http.Cookie{
		Name:     opts.name,
		Value:    base64.URLEncoding.EncodeToString(token),
		Expires:  time.Now().UTC().Add(opts.lifetime),
		HttpOnly: opts.httponly,
		Secure:   opts.secure,
		Path:     "/",
	}

	http.SetCookie(w, &cookie)
}
