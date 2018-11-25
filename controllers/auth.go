package controllers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/alexedwards/scs"
	log "github.com/sirupsen/logrus"

	"github.com/ctf-zone/ctfzone/internal/mailer"
	"github.com/ctf-zone/ctfzone/models"
)

// AuthRegister handles new user registration
func AuthRegister(db *models.Repository, m mailer.Sender) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		var u models.User

		if err := jsonDecode(r, &u); err != nil {
			handleError(w, r, ErrInvalidJSON)
			return
		}

		if valid, ferrs := userValidate(db, &u); !valid {
			handleError(w, r, ErrDuplicate.SetFieldsErrors(ferrs))
			return
		}

		if err := db.UsersInsert(&u); err != nil {
			handleError(w, r, ErrInternal.SetError(err))
			return
		}

		if m != nil {

			t, err := db.TokensNew(u.ID, models.TokenTypeActivate, time.Hour*24)
			if err != nil {
				handleError(w, r, ErrInternal.SetError(err))
				return
			}

			if err := db.TokensInsert(t); err != nil {
				handleError(w, r, ErrInternal.SetError(err))
				return
			}

			go func() {
				data := struct {
					Name  string
					Token string
				}{u.Name, t.Token}

				if err := m.Send("activate", u.Email, data); err != nil {
					log.Error(err)
				}
			}()

		} else {
			u.IsActivated = true

			if err := db.UsersUpdate(&u); err != nil {
				handleError(w, r, ErrInternal.SetError(err))
				return
			}
		}

		log.WithFields(log.Fields{
			"controller": "auth",
			"action":     "register",
			"user":       u.Name,
			"email":      u.Email,
		}).Info("User registered")

		responseCreated(w)
	}
}

// AuthLogin handles user login
func AuthLogin(db *models.Repository, sm *scs.Manager) http.HandlerFunc {

	type Request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		var req Request

		if err := jsonDecode(r, &req); err != nil {
			handleError(w, r, ErrInvalidJSON)
			return
		}

		u, err := db.UsersLogin(req.Email, req.Password)

		if err != nil {
			handleError(w, r, ErrInvalidCreds)
			return
		}

		if !u.IsActivated {
			handleError(w, r, ErrAccountIsNotActivated)
			return
		}

		session := sm.Load(r)

		if err := session.RenewToken(w); err != nil {
			handleError(w, r, ErrInternal.SetError(err))
			return
		}

		if err := session.PutInt64(w, "userId", u.ID); err != nil {
			handleError(w, r, ErrInternal.SetError(err))
			return
		}

		log.WithFields(log.Fields{
			"controller": "auth",
			"action":     "login",
			"user":       u.Name,
			"email":      u.Email,
		}).Info("User logged in")

		responseOK(w)
	}
}

// AuthActivate handles user activation
func AuthActivate(db *models.Repository) http.HandlerFunc {

	type Request struct {
		Token string `json:"token"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		var req Request

		err := jsonDecode(r, &req)
		if err != nil {
			handleError(w, r, ErrInvalidJSON)
			return
		}

		t, err := db.TokensOneByTokenAndType(req.Token, models.TokenTypeActivate)

		if err == sql.ErrNoRows {
			handleError(w, r, ErrTokenNotFound)
			return
		} else if err != nil {
			handleError(w, r, ErrInternal.SetError(err))
			return
		}

		if t.ExpiresAt.Before(time.Now()) {
			handleError(w, r, ErrTokenIsExpired)
			return
		}

		u, err := db.UsersOneByID(t.UserID)
		if err != nil {
			handleError(w, r, ErrInternal.SetError(err))
			return
		}

		u.IsActivated = true

		if err := db.UsersUpdate(u); err != nil {
			handleError(w, r, ErrInternal.SetError(err))
			return
		}

		if err := db.TokensDelete(t.ID); err != nil {
			handleError(w, r, ErrInternal.SetError(err))
			return
		}

		log.WithFields(log.Fields{
			"controller": "auth",
			"action":     "activate",
			"user":       u.Name,
			"email":      u.Email,
		}).Info("User activated")

		responseOK(w)
	}
}

// AuthLogout handles user logout
func AuthLogout(sm *scs.Manager) http.HandlerFunc {

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

// AuthSendToken handles sending tokens to user
func AuthSendToken(db *models.Repository, m mailer.Sender) http.HandlerFunc {

	type Request struct {
		Email string `json:"email"`
		Type  string `json:"type"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		var req Request

		if err := jsonDecode(r, &req); err != nil {
			handleError(w, r, ErrInvalidJSON)
			return
		}

		u, err := db.UsersOneByEmail(req.Email)
		if err == sql.ErrNoRows {
			handleError(w, r, ErrUserNotFound)
			return
		} else if err != nil {
			handleError(w, r, ErrInternal.SetError(err))
			return
		}

		t, err := db.TokensNew(u.ID, models.TokenType(req.Type), time.Hour*24)
		if err != nil {
			handleError(w, r, ErrInternal.SetError(err))
			return
		}

		if err := db.TokensInsert(t); err != nil {
			handleError(w, r, ErrInternal.SetError(err))
			return
		}

		go func() {
			data := struct {
				Name  string
				Token string
			}{u.Name, t.Token}

			if err := m.Send(req.Type, u.Email, data); err != nil {
				log.Error(err)
			}
		}()

		log.WithFields(log.Fields{
			"controller": "auth",
			"action":     "send-token",
			"user":       u.Name,
			"email":      u.Email,
		}).Info("Token sent to user")

		responseOK(w)
	}
}

// AuthResetPassword handles password resetting
func AuthResetPassword(db *models.Repository) http.HandlerFunc {

	type Request struct {
		Token    string `json:"token"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		var req Request

		if err := jsonDecode(r, &req); err != nil {
			handleError(w, r, ErrInvalidJSON)
			return
		}

		t, err := db.TokensOneByTokenAndType(req.Token, models.TokenTypeReset)
		if err == sql.ErrNoRows {
			handleError(w, r, ErrTokenNotFound)
			return
		} else if err != nil {
			handleError(w, r, ErrInternal.SetError(err))
			return
		}

		if t.ExpiresAt.Before(time.Now().UTC()) {
			handleError(w, r, ErrTokenIsExpired)
			return
		}

		u, err := db.UsersOneByID(t.UserID)
		if err != nil {
			handleError(w, r, ErrInternal.SetError(err))
			return
		}

		u.Password = req.Password

		if err := db.UsersUpdate(u); err != nil {
			handleError(w, r, ErrInternal.SetError(err))
			return
		}

		if err := db.TokensDelete(t.ID); err != nil {
			handleError(w, r, ErrInternal.SetError(err))
			return
		}

		log.WithFields(log.Fields{
			"controller": "auth",
			"action":     "reset-password",
			"name":       u.Name,
			"email":      u.Email,
		}).Info("User reset password")

		responseOK(w)
	}
}

func AuthCheck(sm *scs.Manager) http.HandlerFunc {

	type Response struct {
		IsLoggedIn bool `json:"isLoggedIn"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		session := sm.Load(r)

		res := &Response{}

		userID, err := session.GetInt("userId")
		if err != nil || userID <= 0 {
			res.IsLoggedIn = false
		} else {
			res.IsLoggedIn = true
		}

		responseJSON(w, res)
	}
}
