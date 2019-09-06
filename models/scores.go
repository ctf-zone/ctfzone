package models

import (
	"fmt"
	"time"

	"github.com/ctf-zone/ctfzone/config"
)

type scoresListOptions struct {
	Pagination
}

type scoresOption func(*scoresListOptions)

func ScoresPagination(p Pagination) scoresOption {
	return func(params *scoresListOptions) {
		params.Pagination = p
	}
}

// Score represents line in scoreboard.
type Score struct {
	UserP     `json:"user"`
	Score     int        `db:"score"      json:"score"`
	Rank      int        `db:"rank"       json:"rank,omitempty"`
	UpdatedAt *time.Time `db:"updated_at" json:"updatedAt,omitempty"`
}

type Scores []*Score

func (l Scores) First() int64 {
	if len(l) == 0 {
		return 0
	}
	return l[0].ID
}

func (l Scores) Last() int64 {
	if len(l) == 0 {
		return 0
	}
	return l[len(l)-1].ID
}

func (l Scores) Len() int {
	return len(l)
}

func scoresQuery(cfg *config.Scoring) string {
	inner := "SELECT users.id, users.name, users.extra, MAX(solutions.created_at) AS updated_at, "

	switch cfg.Type {
	default:
		fallthrough

	case "classic":
		// Simply calculate sum of points of solved challenges.
		inner += "COALESCE(SUM(challenges.points), 0) AS score "

	case "dynamic":
		p := cfg.Dynamic

		// Calculate sum of solved challenges points.
		// points = min + (max - min) * coeff ^ (n - 1)
		formula := fmt.Sprintf(
			"FLOOR(%d + %d * POWER(%.3f, %s - 1))",
			p.Min,
			p.Max-p.Min,
			p.Coeff,
			"challenges_solutions.count",
		)

		inner += fmt.Sprintf("COALESCE(SUM(%s), 0) AS score ", formula)
	}

	inner += "FROM users " +
		"LEFT JOIN solutions ON solutions.user_id = users.id " +
		"LEFT JOIN challenges ON challenges.id = solutions.challenge_id " +
		"LEFT JOIN challenges_solutions ON challenges_solutions.challenge_id = challenges.id " +
		"GROUP BY users.id, users.name, users.extra " +
		"ORDER BY score DESC, updated_at, name"

	query := fmt.Sprintf("SELECT *, ROW_NUMBER() OVER() AS rank FROM (%s) AS s", inner)

	return query
}

func (r *Repository) ScoresList(cfg *config.Scoring, options ...scoresOption) ([]*Score, *PagesInfo, error) {
	list := make(Scores, 0)

	var params scoresListOptions
	for _, opt := range options {
		opt(&params)
	}

	err := r.db.Select(&list, scoresQuery(cfg))
	if err != nil {
		return nil, nil, err
	}

	return list, nil, nil
}

func (r *Repository) ScoresOneByUserID(cfg *config.Scoring, userID int64) (*Score, error) {
	var o Score

	query := fmt.Sprintf("SELECT * FROM (%s) AS scores WHERE id = $1", scoresQuery(cfg))

	err := r.db.Get(&o, query, userID)
	if err != nil {
		return nil, err
	}

	return &o, nil
}
