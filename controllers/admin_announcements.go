package controllers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/ctf-zone/ctfzone/models"
	"github.com/go-chi/chi"
)

func AdminAnnouncementsCreate(db *models.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var a models.Announcement

		if err := jsonDecode(r, &a); err != nil {
			handleError(w, r, ErrInvalidJSON)
			return
		}

		if err := db.AnnouncementsInsert(&a); err != nil {
			handleError(w, r, ErrInternal.SetError(err))
			return
		}

		if err := responseJSON(w, a); err != nil {
			handleError(w, r, ErrInternal.SetError(err))
			return
		}
	}
}

func AdminAnnouncementsUpdate(db *models.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var a models.Announcement

		if err := jsonDecode(r, &a); err != nil {
			handleError(w, r, ErrInvalidJSON)
			return
		}

		announcementID, err := strconv.ParseInt(chi.URLParam(r, "announcementId"), 10, 32)
		if err != nil {
			handleError(w, r, ErrInvalidID)
			return
		}

		a.ID = announcementID

		err = db.AnnouncementsUpdate(&a)
		if err == sql.ErrNoRows {
			handleError(w, r, ErrAnnouncementNotFound)
			return
		} else if err != nil {
			handleError(w, r, ErrInternal.SetError(err))
			return
		}

		if err := responseJSON(w, a); err != nil {
			handleError(w, r, ErrInternal.SetError(err))
			return
		}
	}
}

func AdminAnnouncementsDelete(db *models.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		announcementID, err := strconv.ParseInt(chi.URLParam(r, "announcementId"), 10, 32)
		if err != nil {
			handleError(w, r, ErrInvalidID)
			return
		}

		err = db.AnnouncementsDelete(announcementID)
		if err == sql.ErrNoRows {
			handleError(w, r, ErrAnnouncementNotFound)
			return
		} else if err != nil {
			handleError(w, r, ErrInternal.SetError(err))
			return
		}

		responseOK(w)
	}
}
