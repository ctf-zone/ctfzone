package models

import (
	"fmt"
	"strings"
	"time"
)

type Announcement struct {
	ID          int64     `db:"id,omitempty" json:"id"`
	Title       string    `db:"title"        json:"title"`
	Body        string    `db:"body"         json:"body"`
	ChallengeID *int64    `db:"challenge_id" json:"challengeId,omitempty"`
	CreatedAt   time.Time `db:"created_at"   json:"createdAt"`
	UpdatedAt   time.Time `db:"updated_at"   json:"updatedAt"`
}

type Announcements []Announcement

func (l Announcements) First() int64 {
	if len(l) == 0 {
		return 0
	}
	return l[0].ID
}

func (l Announcements) Last() int64 {
	if len(l) == 0 {
		return 0
	}
	return l[len(l)-1].ID
}

func (l Announcements) Len() int {
	return len(l)
}

type announcementsOptions struct {
	Pagination
	AnnouncementsFilters
}

type announcementsOption func(*announcementsOptions)

func AnnouncementsListPagination(p Pagination) announcementsOption {
	return func(params *announcementsOptions) {
		params.Pagination = p
	}
}

type AnnouncementsFilters struct {
	Title       string `schema:"title"`
	ChallengeID *int64 `schema:"challengeId"`
}

func AnnouncementsFilter(f AnnouncementsFilters) announcementsOption {
	return func(options *announcementsOptions) {
		options.AnnouncementsFilters = f
	}
}

func AnnouncementsPagination(p Pagination) announcementsOption {
	return func(options *announcementsOptions) {
		if !p.IsZero() {
			options.Pagination = p
		} else {
			options.Pagination = defaultPagination
		}
	}
}

func (r *Repository) AnnouncementsInsert(o *Announcement) error {
	o.CreatedAt = now()
	o.UpdatedAt = o.CreatedAt

	stmt, err := r.db.PrepareNamed(
		"INSERT INTO announcements (title, body, challenge_id, created_at, updated_at) " +
			"VALUES(:title, :body, :challenge_id, :created_at, :updated_at) " +
			"RETURNING id")

	if err != nil {
		return err
	}

	return stmt.QueryRowx(o).Scan(&o.ID)

}

func (r *Repository) AnnouncementsUpdate(o *Announcement) error {
	o.UpdatedAt = now()

	query := "UPDATE announcements SET " +
		"title = :title, " +
		"body = :body, " +
		"challenge_id = :challenge_id, " +
		"updated_at = :updated_at " +
		"WHERE id = :id RETURNING updated_at"

	stmt, err := r.db.PrepareNamed(query)
	if err != nil {
		return err
	}

	return stmt.QueryRowx(o).Scan(&o.UpdatedAt)
}

func (r *Repository) AnnouncementsDelete(id int64) error {
	return r.db.QueryRow("DELETE FROM announcements WHERE id = $1 RETURNING id", id).Scan(&id)
}

func announcementsApplyFilters(query string, f AnnouncementsFilters) (string, map[string]interface{}) {
	cond := make([]string, 0)
	params := make(map[string]interface{})

	if f.Title != "" {
		cond = append(cond, "title ILIKE :title")
		params["title"] = fmt.Sprintf("%%%s%%", f.Title)
	}

	if f.ChallengeID != nil {
		cond = append(cond, "challenge_id = :challenge_id")
		params["challenge_id"] = *f.ChallengeID
	}

	if len(cond) > 0 {
		query += " WHERE " + strings.Join(cond, " AND ")
	}

	return query, params
}

func (r *Repository) AnnouncementsList(opts ...announcementsOption) (Announcements, *PagesInfo, error) {
	list := make(Announcements, 0)

	var options announcementsOptions
	for _, opt := range opts {
		opt(&options)
	}

	query, params := announcementsApplyFilters("SELECT * FROM announcements", options.AnnouncementsFilters)
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

func (r *Repository) AnnouncementsOneByID(id int64) (*Announcement, error) {
	var o Announcement

	err := r.db.Get(&o, "SELECT * FROM announcements WHERE id = $1", id)

	if err != nil {
		return nil, err
	}

	return &o, nil
}
