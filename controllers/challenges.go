package controllers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/alexedwards/scs"
	"github.com/go-chi/chi"

	"github.com/ctf-zone/ctfzone/config"
	"github.com/ctf-zone/ctfzone/models"
)

// ChallengesList handles get all challenges request
func ChallengesList(cfg *config.Scoring, db *models.Repository, sm *scs.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session := sm.Load(r)

		userID, err := session.GetInt64("userId")
		if err != nil {
			handleError(w, r, ErrUnauthorizedRequest)
			return
		}

		challenges, _, err := db.ChallengesListE(
			cfg,
			models.ChallengesIncludeMeta(),
			models.ChallengesIncludeUser(userID),
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

// ChallengesGet handles get challenge request
func ChallengesGet(cfg *config.Scoring, db *models.Repository, sm *scs.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session := sm.Load(r)

		userID, err := session.GetInt64("userId")
		if err != nil {
			handleError(w, r, ErrUnauthorizedRequest)
			return
		}

		challengeID, err := strconv.ParseInt(chi.URLParam(r, "challengeId"), 10, 32)
		if err != nil {
			handleError(w, r, ErrInvalidID)
			return
		}

		c, err := db.ChallengesOneByIDE(
			cfg,
			challengeID,
			models.ChallengesIncludeMeta(),
			models.ChallengesIncludeUser(userID),
		)

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
