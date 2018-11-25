package controllers

import (
	"net/http"

	"github.com/alexedwards/scs"
	"golang.org/x/crypto/bcrypt"

	"github.com/ctf-zone/ctfzone/config"
	"github.com/ctf-zone/ctfzone/models"
)

// AdminAuthLogin handles admin login
// summary: Authorizes admin
// tags: [auth]
func AdminAuthLogin(c *config.Config, db *models.Repository, sm *scs.Manager) http.HandlerFunc {

	type Request struct {
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var req Request

		if err := jsonDecode(r, &req); err != nil {
			handleError(w, r, ErrInvalidJSON)
			return
		}

		if err := bcrypt.CompareHashAndPassword(
			[]byte(c.Admin.Password),
			[]byte(req.Password),
		); err != nil {
			handleError(w, r, ErrInvalidCreds)
			return
		}

		session := sm.Load(r)

		if err := session.RenewToken(w); err != nil {
			handleError(w, r, ErrInternal.SetError(err))
			return
		}

		if err := session.PutBool(w, "isAdmin", true); err != nil {
			handleError(w, r, ErrInternal.SetError(err))
			return
		}

		responseOK(w)
	}
}

// AdminAuthLogout handles admin logout
// summary: Logout
// tags: [auth]
func AdminAuthLogout(sm *scs.Manager) http.HandlerFunc {

	type Request struct{}

	return func(w http.ResponseWriter, r *http.Request) {

		var req Request

		if err := jsonDecode(r, &req); err != nil {
			handleError(w, r, ErrInvalidJSON)
			return
		}

		session := sm.Load(r)

		if err := session.Clear(w); err != nil {
			handleError(w, r, ErrInternal.SetError(err))
			return
		}

		responseOK(w)
	}
}

// AdminAuthCheck handles auth check request
// description: Returns 200 if admin is logged in
// tags: [auth]
func AdminAuthCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		responseOK(w)
	}
}
