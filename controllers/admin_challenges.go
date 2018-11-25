package controllers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"

	"github.com/ctf-zone/ctfzone/config"
	"github.com/ctf-zone/ctfzone/models"
)

func AdminChallengesList(cfg *config.Scoring, db *models.Repository) http.HandlerFunc {

	type Params struct {
		models.Pagination
		models.ChallengesFilters
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var params Params

		if err := decoder.Decode(&params, r.URL.Query()); err != nil {
			handleError(w, r, ErrInvalidQueryParams.SetMessage(err.Error()))
			return
		}

		challenges, _, err := db.ChallengesListE(
			cfg,
			models.ChallengesIncludeLocked(),
			models.ChallengesIncludeMeta(),
			models.ChallengesFilter(params.ChallengesFilters),
		)
		if err != nil {
			handleError(w, r, ErrInternal.SetError(err))
			return
		}

		if err := responseJSON(w, challenges); err != nil {
			handleError(w, r, ErrInternal.SetError(err))
			return
		}
	}
}

func AdminChallengesGet(cfg *config.Scoring, db *models.Repository) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		challengeID, err := strconv.ParseInt(chi.URLParam(r, "challengeId"), 10, 32)
		if err != nil {
			handleError(w, r, ErrInvalidID)
			return
		}

		challenge, err := db.ChallengesOneByIDE(
			cfg,
			challengeID,
			models.ChallengesIncludeLocked(),
		)
		if err == sql.ErrNoRows {
			handleError(w, r, ErrChallengeNotFound)
			return
		} else if err != nil {
			handleError(w, r, ErrInternal.SetError(err))
			return
		}

		if err := responseJSON(w, challenge); err != nil {
			handleError(w, r, ErrInternal.SetError(err))
			return
		}
	}
}

func AdminChallengesCreate(db *models.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var c models.Challenge

		if err := jsonDecode(r, &c); err != nil {
			handleError(w, r, ErrInvalidJSON)
			return
		}

		if err := db.ChallengesInsert(&c); err != nil {
			handleError(w, r, ErrInternal.SetError(err))
			return
		}

		if err := responseJSON(w, c); err != nil {
			handleError(w, r, ErrInternal.SetError(err))
			return
		}
	}
}

func AdminChallengesUpdate(db *models.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var c models.Challenge

		if err := jsonDecode(r, &c); err != nil {
			handleError(w, r, ErrInvalidJSON)
			return
		}

		// description: Numeric ID of the challenge
		// schema: { type: integer }
		challengeID, err := strconv.ParseInt(chi.URLParam(r, "challengeId"), 10, 32)
		if err != nil {
			handleError(w, r, ErrInvalidID)
			return
		}

		c.ID = challengeID

		err = db.ChallengesUpdate(&c)
		if err == sql.ErrNoRows {
			handleError(w, r, ErrChallengeNotFound)
			return
		} else if err != nil {
			handleError(w, r, ErrInternal.SetError(err))
			return
		}

		if err := responseJSON(w, c); err != nil {
			handleError(w, r, ErrInternal.SetError(err))
			return
		}
	}
}

func AdminChallengesDelete(db *models.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// description: Numeric ID of the challenge
		// schema: { type: integer }
		challengeID, err := strconv.ParseInt(chi.URLParam(r, "challengeId"), 10, 32)
		if err != nil {
			handleError(w, r, ErrInvalidID)
			return
		}

		err = db.ChallengesDelete(challengeID)
		if err == sql.ErrNoRows {
			handleError(w, r, ErrChallengeNotFound)
			return
		} else if err != nil {
			handleError(w, r, ErrInternal.SetError(err))
			return
		}

		responseOK(w)
	}
}
