package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ctf-zone/ctfzone/config"
	"github.com/ctf-zone/ctfzone/models"
)

// ScoresList handles scores request
// summary: Returns list of scores
// tags: [scores]
func ScoresList(cfg *config.Scoring, db *models.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		scores, _, err := db.ScoresList(cfg)

		if err != nil {
			handleError(w, r, ErrInternal.SetError(err))
			return
		}

		// schema: Scores
		if err := responseJSON(w, scores); err != nil {
			handleError(w, r, ErrInternal.SetError(err))
			return
		}
	}
}

// ScoresCtftimeList returns scoreboard in ctftime format
// summary: Returns scores in CTFTime format
// description: https://ctftime.org/json-scoreboard-feed
// tags: [scores]
func ScoresCtftimeList(cfg *config.Scoring, db *models.Repository) http.HandlerFunc {

	type ScoreCtftime struct {
		Pos        int             `json:"pos"`
		Team       json.RawMessage `json:"team"`
		Score      int             `json:"score"`
		LastAccept int64           `json:"lastAccept"`
	}

	type Ctftime struct {
		Standings []*ScoreCtftime `json:"standings"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		scores, _, err := db.ScoresList(cfg)

		if err != nil {
			handleError(w, r, ErrInternal.SetError(err))
			return
		}

		ctftime := &Ctftime{}

		standings := make([]*ScoreCtftime, 0)

		for _, s := range scores {
			if s.Score != 0 {
				standings = append(standings, &ScoreCtftime{
					Pos:        s.Rank,
					Score:      s.Score,
					Team:       []byte(fmt.Sprintf("%+q", s.UserP.Name)),
					LastAccept: s.UpdatedAt.Unix(),
				})
			}
		}

		ctftime.Standings = standings

		// schema: ScoresCTFTime
		if err := responseJSON(w, ctftime); err != nil {
			handleError(w, r, ErrInternal.SetError(err))
			return
		}
	}
}
