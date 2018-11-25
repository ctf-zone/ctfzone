package controllers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"

	"github.com/ctf-zone/ctfzone/config"
	"github.com/ctf-zone/ctfzone/models"
)

// AdminUsersList returns list of registered users
func AdminUsersList(cfg *config.Config, db *models.Repository) http.HandlerFunc {

	type Params struct {
		models.Pagination
		models.UsersFilters
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var params Params

		if err := decoder.Decode(&params, r.URL.Query()); err != nil {
			handleError(w, r, ErrInvalidQueryParams.SetMessage(err.Error()))
			return
		}

		users, pagesInfo, err := db.UsersList(
			models.UsersPagination(params.Pagination),
			models.UsersFilter(params.UsersFilters),
		)
		if err != nil {
			handleError(w, r, ErrInternal.SetError(err))
			return
		}

		if pagesInfo != nil {
			addLinkHeader(w, cfg.Server.AdminBaseURL()+r.URL.Path, pagesInfo, r.URL.Query())
		}

		if err := responseJSON(w, users); err != nil {
			handleError(w, r, ErrInternal.SetError(err))
			return
		}
	}
}

// AdminUsersGet handler get user by id request
func AdminUsersGet(db *models.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userID, err := strconv.ParseInt(chi.URLParam(r, "userId"), 10, 32)
		if err != nil {
			handleError(w, r, ErrInvalidID)
			return
		}

		u, err := db.UsersOneByID(userID)
		if err == sql.ErrNoRows {
			handleError(w, r, ErrUserNotFound)
			return
		} else if err != nil {
			handleError(w, r, ErrInternal.SetError(err))
			return
		}

		if err := responseJSON(w, u); err != nil {
			handleError(w, r, ErrInternal.SetError(err))
			return
		}
	}
}

// AdminUsersCreate handles user create request
func AdminUsersCreate(db *models.Repository) http.HandlerFunc {
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

		if err := responseJSON(w, u); err != nil {
			handleError(w, r, ErrInternal.SetError(err))
			return
		}
	}
}

// AdminUsersList handles user update request
func AdminUsersUpdate(db *models.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u models.User

		if err := jsonDecode(r, &u); err != nil {
			handleError(w, r, ErrInvalidJSON)
			return
		}

		userID, err := strconv.ParseInt(chi.URLParam(r, "userId"), 10, 32)
		if err != nil {
			handleError(w, r, ErrInvalidID)
			return
		}

		u.ID = userID

		err = db.UsersUpdate(&u)
		if err == sql.ErrNoRows {
			handleError(w, r, ErrUserNotFound)
			return
		} else if err != nil {
			handleError(w, r, ErrInternal.SetError(err))
			return
		}

		if err := responseJSON(w, u); err != nil {
			handleError(w, r, ErrInternal.SetError(err))
			return
		}
	}
}

// AdminUsersList handles user delete request
func AdminUsersDelete(db *models.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userID, err := strconv.ParseInt(chi.URLParam(r, "userId"), 10, 32)
		if err != nil {
			handleError(w, r, ErrInvalidID)
			return
		}

		err = db.UsersDelete(userID)
		if err == sql.ErrNoRows {
			handleError(w, r, ErrUserNotFound)
			return
		} else if err != nil {
			handleError(w, r, ErrInternal.SetError(err))
			return
		}

		responseOK(w)
	}
}
