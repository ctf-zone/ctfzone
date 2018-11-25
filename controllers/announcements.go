package controllers

import (
	"net/http"
	"strconv"

	"github.com/ctf-zone/ctfzone/models"
	"github.com/go-chi/chi"
)

// AnnouncementsList returns all announcements
func AnnouncementsList(db *models.Repository) http.HandlerFunc {

	type Params struct {
		models.AnnouncementsFilters
		models.Pagination
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var params Params

		if err := decoder.Decode(&params, r.URL.Query()); err != nil {
			handleError(w, r, ErrInvalidQueryParams.SetMessage(err.Error()))
			return
		}

		announcements, _, err := db.AnnouncementsList(
			models.AnnouncementsFilter(params.AnnouncementsFilters),
			models.AnnouncementsPagination(params.Pagination),
		)
		if err != nil {
			handleError(w, r, ErrInternal.SetError(err))
			return
		}

		// schema: Announcements
		if err := responseJSON(w, announcements); err != nil {
			handleError(w, r, ErrInternal.SetError(err))
			return
		}
	}
}

func AnnouncementsGet(db *models.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		announcementID, err := strconv.ParseInt(chi.URLParam(r, "announcementId"), 10, 32)
		if err != nil {
			handleError(w, r, ErrInvalidID)
			return
		}

		announcement, err := db.AnnouncementsOneByID(announcementID)
		if err != nil {
			handleError(w, r, ErrInternal.SetError(err))
			return
		}

		if err := responseJSON(w, announcement); err != nil {
			handleError(w, r, ErrInternal.SetError(err))
			return
		}
	}
}
