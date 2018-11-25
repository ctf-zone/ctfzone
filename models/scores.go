package models

import (
	"fmt"
	"time"

	"github.com/ctf-zone/ctfzone/config"
	udb "upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

type ScoresRepository interface {
	List(...scoresOption) ([]*Score, *PagesInfo, error)
}

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

func (r *Repository) scoresQuery(cfg *config.Scoring, options scoresListOptions) sqlbuilder.Selector {
	scoresInner := r.db.
		SelectFrom("users").
		Columns(
			"users.id",
			"users.name",
			"users.extra",
			udb.Raw("MAX(solutions.created_at) AS updated_at"),
		).
		LeftJoin("solutions").On("solutions.user_id = users.id").
		LeftJoin("challenges").On("challenges.id = solutions.challenge_id").
		GroupBy(
			"users.id",
			"users.name",
			"users.extra",
		).
		OrderBy(
			"-score",
			"updated_at",
			"name",
		)

	switch cfg.Type {

	default:
		fallthrough

	case "classic":
		// Simply calculate sum of points of solved challenges.
		scoresInner = scoresInner.
			Columns(
				udb.Raw("COALESCE(SUM(challenges.points), 0) AS score"),
			)

	case "dynamic":
		p := cfg.Dynamic

		// Calculate sum of solved challenges points.
		// points = min + (max - min) * coeff ^ (n - 1)
		formula := fmt.Sprintf(
			"FLOOR(%d + %d * POWER(%.2f, %s - 1))",
			p.Min,
			p.Max-p.Min,
			p.Coeff,
			"challenges_solutions.count",
		)

		scoresInner = scoresInner.
			Columns(
				udb.Raw(fmt.Sprintf("COALESCE(SUM(%s), 0) AS score", formula)),
			).
			LeftJoin("challenges_solutions").On("challenges_solutions.challenge_id = challenges.id")
	}

	query := r.db.
		SelectFrom(scoresInner).
		As("scores").
		Columns(
			"*",
			udb.Raw("ROW_NUMBER() OVER() AS rank"),
		)

	return query
}

func (r *Repository) ScoresList(cfg *config.Scoring, options ...scoresOption) ([]*Score, *PagesInfo, error) {
	list := make(Scores, 0)

	var params scoresListOptions
	for _, opt := range options {
		opt(&params)
	}

	query := r.db.
		SelectFrom(r.scoresQuery(cfg, params)).
		As("scores").
		Paginate(params.Pagination.Count).
		Cursor("rank")

	if err := paginate(query, params.Pagination).All(&list); err != nil {
		return nil, nil, err
	}

	pi, err := r.pagesInfo(query, params.Pagination, list)

	if err != nil {
		return nil, nil, err
	}

	return list, pi, nil
}

func (r *Repository) ScoresOneByUserID(cfg *config.Scoring, userID int64) (*Score, error) {
	var o Score

	err := r.db.
		SelectFrom(r.scoresQuery(cfg, scoresListOptions{})).
		As("scores").
		Where("id", userID).
		Limit(1).
		One(&o)

	if err != nil {
		return nil, handleErr(err)
	}

	return &o, nil
}
