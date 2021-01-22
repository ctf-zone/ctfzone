package controllers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/alexedwards/scs"
	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"

	"github.com/ctf-zone/ctfzone/config"
	"github.com/ctf-zone/ctfzone/internal/crypto"
	"github.com/ctf-zone/ctfzone/models"
)

// UserSolutionsCreate handles challenge solution request
func UserSolutionsCreate(cfg *config.Scoring, db *models.Repository, sm *scs.Manager) http.HandlerFunc {

	type Request struct {
		Flag string `json:"flag"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		var req Request

		if err := jsonDecode(r, &req); err != nil {
			handleError(w, r, ErrInvalidJSON)
			return
		}

		session := sm.Load(r)

		userID, err := session.GetInt64("userId")
		if err != nil {
			handleError(w, r, ErrUnauthorizedRequest)
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

		challengeID, err := strconv.ParseInt(chi.URLParam(r, "challengeId"), 10, 64)
		if err != nil {
			handleError(w, r, ErrInvalidID)
			return
		}

		t, err := db.ChallengesOneByIDE(cfg, challengeID, models.ChallengesIncludeUser(userID))
		if err == sql.ErrNoRows {
			handleError(w, r, ErrChallengeNotFound)
			return
		} else if err != nil {
			handleError(w, r, ErrInternal.SetError(err))
			return
		}

		s, _ := db.SolutionsOneByID(userID, challengeID)
		if s != nil {
			handleError(w, r, ErrDuplicate)
			return
		}

		if err := crypto.CheckFlag(t.FlagHash, req.Flag); err != nil {

			// Log attempts.
			log.WithFields(log.Fields{
				"controller": "challenge",
				"action":     "solve",
				"username":   u.Name,
				"email":      u.Email,
				"challenge":  t.Title,
				"flag":       req.Flag,
			}).Info("Challenge solve")

			handleError(w, r, ErrInvalidFlag)
			return
		}

		s = &models.Solution{
			UserID:      userID,
			ChallengeID: challengeID,
		}

		if err := db.SolutionsInsert(s); err != nil {
			handleError(w, r, ErrInternal.SetError(err))
			return
		}

		log.WithFields(log.Fields{
			"controller": "challenge",
			"action":     "solve",
			"username":   u.Name,
			"email":      u.Email,
			"challenge":  t.Title,
		}).Info("Challenge solve")

		responseOK(w)
	}
}

// UserLikesCreate handles like creation for challenge
func UserLikesCreate(db *models.Repository, sm *scs.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session := sm.Load(r)

		userID, err := session.GetInt64("userId")
		if err != nil {
			handleError(w, r, ErrUnauthorizedRequest)
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

		challengeID, err := strconv.ParseInt(chi.URLParam(r, "challengeId"), 10, 64)
		if err != nil {
			handleError(w, r, ErrInvalidID)
			return
		}

		t, err := db.ChallengesOneByID(challengeID)
		if err == sql.ErrNoRows {
			handleError(w, r, ErrChallengeNotFound)
			return
		} else if err != nil {
			handleError(w, r, ErrInternal.SetError(err))
			return
		}

		if l, _ := db.LikesOneByID(userID, challengeID); l != nil {
			handleError(w, r, ErrDuplicate)
			return
		}

		l := &models.Like{
			ChallengeID: challengeID,
			UserID:      userID,
		}

		if err := db.LikesInsert(l); err != nil {
			handleError(w, r, ErrInternal.SetError(err))
			return
		}

		log.WithFields(log.Fields{
			"controller": "challenge",
			"action":     "like-create",
			"username":   u.Name,
			"email":      u.Email,
			"challenge":  t.Title,
		}).Info("Challenge liked")

		responseCreated(w)
	}
}

// UserLikesDelete handles like deletion
func UserLikesDelete(db *models.Repository, sm *scs.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session := sm.Load(r)

		userID, err := session.GetInt64("userId")
		if err != nil {
			handleError(w, r, ErrUnauthorizedRequest)
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

		challengeID, err := strconv.ParseInt(chi.URLParam(r, "challengeId"), 10, 64)
		if err != nil {
			handleError(w, r, ErrInvalidID)
			return
		}

		t, err := db.ChallengesOneByID(challengeID)
		if err == sql.ErrNoRows {
			handleError(w, r, ErrChallengeNotFound)
			return
		} else if err != nil {
			handleError(w, r, ErrInternal.SetError(err))
			return
		}

		err = db.LikesDelete(userID, challengeID)

		if err != nil && err != sql.ErrNoRows {
			handleError(w, r, ErrInternal.SetError(err))
			return
		}

		log.WithFields(log.Fields{
			"controller": "user",
			"action":     "like-delete",
			"username":   u.Name,
			"email":      u.Email,
			"challenge":  t.Title,
		}).Info("Challenge like reset")

		responseOK(w)
	}
}

func UserGetStats(cfg *config.Config, db *models.Repository, sm *scs.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session := sm.Load(r)

		userID, err := session.GetInt64("userId")
		if err != nil {
			handleError(w, r, ErrUnauthorizedRequest)
			return
		}

		s, err := db.ScoresOneByUserID(&cfg.Game.Scoring, userID)

		if err != nil {
			handleError(w, r, ErrInternal.SetError(err))
			return
		}

		responseJSON(w, s)
	}
}
