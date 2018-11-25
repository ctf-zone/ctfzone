package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
	udb "upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
	"upper.io/db.v3/postgresql"

	"github.com/ctf-zone/ctfzone/config"
	"github.com/ctf-zone/ctfzone/internal/crypto"
)

type Difficulty string

const (
	DifficultyEasy   = Difficulty("easy")
	DifficultyMedium = Difficulty("medium")
	DifficultyHard   = Difficulty("hard")
)

type Challenge struct {
	ID          int64                  `db:"id,omitempty"         json:"id"`
	Title       string                 `db:"title"                json:"title"`
	Categories  postgresql.StringArray `db:"categories,omitempty" json:"categories"`
	Points      int                    `db:"points"               json:"points"`
	Description string                 `db:"description"          json:"description"`
	Difficulty  Difficulty             `db:"difficulty"           json:"difficulty"`
	Flag        string                 `db:"-"                    json:"flag,omitempty"`
	FlagHash    string                 `db:"flag_hash"            json:"-"`
	IsLocked    bool                   `db:"is_locked"            json:"isLocked"`
	CreatedAt   time.Time              `db:"created_at"           json:"createdAt"`
	UpdatedAt   time.Time              `db:"updated_at"           json:"updatedAt"`
}

type ChallengeMeta struct {
	SolutionsCount int `db:"solutions_count" json:"solutionsCount"`
	LikesCount     int `db:"likes_count"     json:"likesCount"`
	HintsCount     int `db:"hints_count"     json:"hintsCount"`
}

type ChallengeUser struct {
	IsSolved bool `db:"is_solved" json:"isSolved"`
	IsLiked  bool `db:"is_liked"  json:"isLiked"`
}

type ChallengeE struct {
	Challenge Challenge      `json:"challenge"`
	Meta      *ChallengeMeta `json:"meta,omitempty"`
	User      *ChallengeUser `json:"user,omitempty"`
}

type ChallengesE []ChallengeE

func (l ChallengesE) First() int64 {
	if len(l) == 0 {
		return 0
	}
	return l[0].Challenge.ID
}

func (l ChallengesE) Last() int64 {
	if len(l) == 0 {
		return 0
	}
	return l[len(l)-1].Challenge.ID
}

func (l ChallengesE) Len() int {
	return len(l)
}

type challengesListOptions struct {
	Pagination
	ChallengesFilters
	IncludeMeta   bool
	UserID        int64
	IncludeLocked bool
}

type challengesOption func(*challengesListOptions)

func ChallengesPagination(p Pagination) challengesOption {
	return func(params *challengesListOptions) {
		if !p.IsZero() {
			params.Pagination = p
		} else {
			params.Pagination = defaultPagination
		}
	}
}

func ChallengesIncludeMeta() challengesOption {
	return func(params *challengesListOptions) {
		params.IncludeMeta = true
	}
}

func ChallengesIncludeUser(userID int64) challengesOption {
	return func(params *challengesListOptions) {
		params.UserID = userID
	}
}

func ChallengesIncludeLocked() challengesOption {
	return func(params *challengesListOptions) {
		params.IncludeLocked = true
	}
}

type ChallengesFilters struct {
	Title      string   `schema:"title"`
	Categories []string `schema:"categories"`
	IsLocked   *bool    `schema:"isLocked"`
}

func ChallengesFilter(f ChallengesFilters) challengesOption {
	return func(params *challengesListOptions) {
		params.ChallengesFilters = f
	}
}

func (r *Repository) ChallengesInsert(o *Challenge) error {

	o.CreatedAt = now()
	o.UpdatedAt = o.CreatedAt
	o.FlagHash = crypto.HashFlag(o.Flag)
	o.Flag = ""

	row, err := r.db.
		InsertInto("challenges").
		Values(o).
		Returning("id").
		QueryRow()

	if err != nil {
		return err
	}

	return row.Scan(&o.ID)
}

func (r *Repository) ChallengesUpdate(o *Challenge) error {
	o.UpdatedAt = now()

	query := r.db.
		Update("challenges").
		Set(map[string]interface{}{
			"title":       o.Title,
			"categories":  o.Categories,
			"points":      o.Points,
			"description": o.Description,
			"difficulty":  o.Difficulty,
			"is_locked":   o.IsLocked,
			"updated_at":  o.UpdatedAt,
		}).
		Where("id", o.ID)

	if o.Flag != "" {
		o.FlagHash = crypto.HashFlag(o.Flag)
		o.Flag = ""
		query = query.Set("flag_hash", o.FlagHash)
	}

	res, err := query.Exec()

	if err != nil {
		return err
	}

	if n, err := res.RowsAffected(); err != nil {
		return err
	} else if n != 1 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *Repository) ChallengesDelete(id int64) error {
	res, err := r.db.
		DeleteFrom("challenges").
		Where("id", id).
		Exec()

	if err != nil {
		return err
	}

	if n, err := res.RowsAffected(); err != nil {
		return err
	} else if n != 1 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *Repository) challengesQuery(cfg *config.Scoring, params challengesListOptions) sqlbuilder.Selector {
	query := r.db.
		SelectFrom("challenges").
		Columns(
			"challenges.id",
			"challenges.title",
			"challenges.categories",
			"challenges.description",
			"challenges.difficulty",
			"challenges.is_locked",
			"challenges.created_at",
			"challenges.updated_at",
		).
		LeftJoin("challenges_solutions").On("challenges_solutions.challenge_id = challenges.id").
		LeftJoin("challenges_likes").On("challenges_likes.challenge_id = challenges.id").
		LeftJoin("challenges_hints").On("challenges_hints.challenge_id = challenges.id")

	// Exclude locked by default to prevent
	// showing them accidently to users.
	if !params.IncludeLocked {
		query = query.Where("is_locked", false)
	}

	if params.IncludeMeta {
		query = query.
			Columns(
				udb.Raw("COALESCE(challenges_solutions.count, 0) AS solutions_count"),
				udb.Raw("COALESCE(challenges_likes.count, 0) AS likes_count"),
				udb.Raw("COALESCE(challenges_hints.count, 0) AS hints_count"),
			)
	}

	if params.UserID != 0 {
		query = query.
			Columns(
				udb.Raw("COALESCE(? = ANY (challenges_solutions.users), FALSE) AS is_solved", params.UserID),
				udb.Raw("COALESCE(? = ANY (challenges_likes.users), FALSE) AS is_liked", params.UserID),
			)
	}

	switch cfg.Type {

	default:
		fallthrough

	case "classic":
		query = query.Columns("challenges.points")

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

		query = query.Columns(
			udb.Raw(fmt.Sprintf("COALESCE(%s, %d) AS points", formula, p.Max)),
		)
	}

	return query
}

func (r *Repository) ChallengesOneByID(id int64) (*Challenge, error) {
	var o Challenge

	err := r.db.
		SelectFrom("challenges").
		Where("id", id).
		Limit(1).
		One(&o)

	if err != nil {
		return nil, handleErr(err)
	}

	return &o, nil
}

func challengesApplyFilters(query sqlbuilder.Selector, f ChallengesFilters) sqlbuilder.Selector {
	if f.Title != "" {
		query = query.Where(
			udb.Cond{
				"title ILIKE": fmt.Sprintf("%%%s%%", f.Title),
			},
		)
	}

	if len(f.Categories) > 0 {
		query = query.Where(
			udb.Cond{
				"categories &&": pq.StringArray(f.Categories),
			},
		)
	}

	if f.IsLocked != nil {
		query = query.Where("is_locked", *f.IsLocked)
	}

	return query
}

func (r *Repository) ChallengesOneByIDE(cfg *config.Scoring, id int64, options ...challengesOption) (*ChallengeE, error) {
	var o ChallengeE

	var params challengesListOptions
	for _, opt := range options {
		opt(&params)
	}

	err := r.db.
		SelectFrom(r.challengesQuery(cfg, params)).
		As("challenges_e").
		Where("id", id).
		Limit(1).
		One(&o)

	if err != nil {
		return nil, handleErr(err)
	}

	return &o, nil
}

func (r *Repository) ChallengesListE(cfg *config.Scoring, options ...challengesOption) ([]ChallengeE, *PagesInfo, error) {
	list := make(ChallengesE, 0)

	var params challengesListOptions

	for _, opt := range options {
		opt(&params)
	}

	query := challengesApplyFilters(
		r.db.SelectFrom(r.challengesQuery(cfg, params)).As("challenges_e"),
		params.ChallengesFilters,
	).
		Paginate(params.Pagination.Count).
		Cursor("-id")

	if err := paginate(query, params.Pagination).All(&list); err != nil {
		return nil, nil, err
	}

	pi, err := r.pagesInfo(query, params.Pagination, list)

	if err != nil {
		return nil, nil, err
	}

	return list, pi, nil
}
