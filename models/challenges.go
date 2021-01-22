package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/lib/pq"

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
	ID          int64          `db:"id,omitempty"         json:"id"`
	Title       string         `db:"title"                json:"title"`
	Categories  pq.StringArray `db:"categories,omitempty" json:"categories"`
	Points      int            `db:"points"               json:"points"`
	Description string         `db:"description"          json:"description"`
	Difficulty  Difficulty     `db:"difficulty"           json:"difficulty"`
	Flag        string         `db:"-"                    json:"flag,omitempty"`
	FlagHash    string         `db:"flag_hash"            json:"-"`
	IsLocked    bool           `db:"is_locked"            json:"isLocked"`
	IsAvailable bool           `db:"is_available"         json:"-"`
	DependsOn   *int64         `db:"depends_on"           json:"dependsOn,omitempty"`
	CreatedAt   time.Time      `db:"created_at"           json:"createdAt"`
	UpdatedAt   time.Time      `db:"updated_at"           json:"updatedAt"`
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
	Challenge      `json:"challenge"`
	*ChallengeMeta `json:"meta,omitempty"`
	*ChallengeUser `json:"user,omitempty"`
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

	fields := []string{
		"title",
		"categories",
		"points",
		"description",
		"difficulty",
		"flag_hash",
		"is_locked",
		"created_at",
		"updated_at",
	}

	if o.DependsOn != nil {
		fields = append(fields, "depends_on")
	}

	placeholders := make([]string, 0)

	for _, f := range fields {
		placeholders = append(placeholders, ":"+f)
	}

	stmt, err := r.db.PrepareNamed(
		"INSERT INTO challenges (" + strings.Join(fields, ", ") + ") " +
			"VALUES(" + strings.Join(placeholders, ", ") + ") " +
			"RETURNING id")

	if err != nil {
		return err
	}

	return stmt.QueryRowx(o).Scan(&o.ID)
}

func (r *Repository) ChallengesUpdate(o *Challenge) error {
	o.UpdatedAt = now()

	query := "UPDATE challenges SET " +
		"title = :title, " +
		"categories = :categories, " +
		"points = :points, " +
		"description = :description, " +
		"difficulty = :difficulty, " +
		"is_locked = :is_locked, " +
		"updated_at = :updated_at"

	if o.DependsOn != nil {
		query += ", depends_on = :depends_on "
	}

	if o.Flag != "" {
		o.FlagHash = crypto.HashFlag(o.Flag)
		o.Flag = ""
		query += ", flag_hash = :flag_hash "
	}

	query += " WHERE id = :id RETURNING updated_at"

	stmt, err := r.db.PrepareNamed(query)
	if err != nil {
		return err
	}

	return stmt.QueryRowx(o).Scan(&o.UpdatedAt)
}

func (r *Repository) ChallengesDelete(id int64) error {
	return r.db.QueryRow("DELETE FROM challenges WHERE id = $1 RETURNING id", id).Scan(&id)
}

func challengesQuery(cfg *config.Scoring, options challengesListOptions) (string, map[string]interface{}) {
	params := make(map[string]interface{})

	query := "SELECT " +
		"challenges.id, " +
		"challenges.title, " +
		"challenges.categories, " +
		"challenges.description, " +
		"challenges.difficulty, " +
		"challenges.is_locked, " +
		"challenges.depends_on, " +
		"challenges.flag_hash, " +
		"challenges.created_at, " +
		"challenges.updated_at, "

	// Scoring
	switch cfg.Type {
	default:
		fallthrough

	case "classic":
		query += "challenges.points"

	case "dynamic":
		p := cfg.Dynamic

		// TODO: cfg.Dynamic method
		// Calculate sum of solved challenges points.
		// points = min + (max - min) * coeff ^ (n - 1)
		formula := fmt.Sprintf(
			"FLOOR(%d + %d * POWER(%.3f, %s - 1))",
			p.Min,
			p.Max-p.Min,
			p.Coeff,
			"challenges_solutions.count",
		)

		query += fmt.Sprintf("COALESCE(%s, %d) AS points", formula, p.Max)
	}

	if options.IncludeMeta {
		query += ", " +
			"COALESCE(challenges_solutions.count, 0) AS solutions_count, " +
			"COALESCE(challenges_likes.count, 0) AS likes_count, " +
			"COALESCE(challenges_hints.count, 0) AS hints_count"
	}

	if options.UserID != 0 {
		query += ", " +
			"COALESCE(:user_id = ANY (challenges_solutions.users), FALSE) AS is_solved, " +
			"COALESCE(:user_id = ANY (challenges_likes.users), FALSE) AS is_liked, " +
			"(NOT challenges.is_locked AND CASE WHEN challenges.depends_on IS NOT NULL THEN COALESCE(:user_id = ANY (dependencies.users), FALSE) ELSE TRUE END) AS is_available "

		params["user_id"] = options.UserID
	}

	query += " FROM challenges " +
		"LEFT JOIN challenges_solutions ON challenges_solutions.challenge_id = challenges.id " +
		"LEFT JOIN challenges_solutions as dependencies ON dependencies.challenge_id = challenges.depends_on " +
		"LEFT JOIN challenges_likes ON challenges_likes.challenge_id = challenges.id " +
		"LEFT JOIN challenges_hints ON challenges_hints.challenge_id = challenges.id "

	// Exclude locked by default to prevent
	// showing them accidently to users.
	if !options.IncludeLocked {
		query = fmt.Sprintf("SELECT * FROM (%s) as q WHERE is_available = :is_available OR is_solved = TRUE", query)
		params["is_available"] = true
	}

	return query, params
}

func (r *Repository) ChallengesOneByID(id int64) (*Challenge, error) {
	var o Challenge

	err := r.db.Get(&o, "SELECT * FROM challenges WHERE id = $1", id)

	if err != nil {
		return nil, err
	}

	return &o, nil
}

func challengesApplyFilters(query string, params map[string]interface{}, f ChallengesFilters) (string, map[string]interface{}) {
	cond := make([]string, 0)

	if f.Title != "" {
		cond = append(cond, "title ILIKE :title")
		params["title"] = fmt.Sprintf("%%%s%%", f.Title)
	}

	if len(f.Categories) > 0 {
		cond = append(cond, "categories && :categories")
		params["categories"] = pq.StringArray(f.Categories)
	}

	if f.IsLocked != nil {
		// Could be added before by "IncludeLocked" option.
		if !strings.Contains(query, "is_locked = :is_locked") {
			cond = append(cond, "is_locked = :is_locked")
		}
		params["is_locked"] = *f.IsLocked
	}

	if len(cond) > 0 {
		query += condPrefix(query) + " " + strings.Join(cond, " AND ")
	}

	return query, params
}

func (r *Repository) ChallengesOneByIDE(cfg *config.Scoring, id int64, opts ...challengesOption) (*ChallengeE, error) {
	var o ChallengeE

	var options challengesListOptions
	for _, opt := range opts {
		opt(&options)
	}

	query, params := challengesQuery(cfg, options)
	query = fmt.Sprintf("SELECT * FROM (%s) AS c WHERE id = :id", query)
	params["id"] = id

	stmt, err := r.db.PrepareNamed(query)
	if err != nil {
		return nil, err
	}

	if err := stmt.Get(&o, params); err != nil {
		return nil, err
	}

	return &o, nil
}

func (r *Repository) ChallengesListE(cfg *config.Scoring, opts ...challengesOption) ([]ChallengeE, *PagesInfo, error) {
	list := make(ChallengesE, 0)

	var options challengesListOptions

	for _, opt := range opts {
		opt(&options)
	}

	query, params := challengesQuery(cfg, options)
	query, params = challengesApplyFilters(query, params, options.ChallengesFilters)
	pageQuery, params := paginate(query, params, options.Pagination)

	stmt, err := r.db.PrepareNamed(pageQuery)
	if err != nil {
		return nil, nil, err
	}

	if err := stmt.Select(&list, params); err != nil {
		return nil, nil, err
	}

	pi, err := r.pagesInfo(query, params, options.Pagination, list)

	if err != nil {
		return nil, nil, err
	}

	return list, pi, nil
}
