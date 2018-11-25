package models

import (
	"database/sql"
	"fmt"
	"time"

	udb "upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
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

	row, err := r.db.
		InsertInto("announcements").
		Values(o).
		Returning("id").
		QueryRow()

	if err != nil {
		return err
	}

	return row.Scan(&o.ID)
}

func (r *Repository) AnnouncementsUpdate(o *Announcement) error {
	o.UpdatedAt = now()

	res, err := r.db.
		Update("announcements").
		Set(map[string]interface{}{
			"title":        o.Title,
			"body":         o.Body,
			"challenge_id": o.ChallengeID,
			"updated_at":   o.UpdatedAt,
		}).
		Where("id", o.ID).
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

func (r *Repository) AnnouncementsDelete(id int64) error {
	res, err := r.db.
		DeleteFrom("announcements").
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

func announcementsApplyFilters(query sqlbuilder.Selector, f AnnouncementsFilters) sqlbuilder.Selector {
	if f.Title != "" {
		query = query.Where(
			udb.Cond{
				"title ILIKE": fmt.Sprintf("%%%s%%", f.Title),
			},
		)
	}

	if f.ChallengeID != nil {
		query = query.Where("challenge_id", *f.ChallengeID)
	}

	return query
}

func (r *Repository) AnnouncementsList(options ...announcementsOption) (Announcements, *PagesInfo, error) {
	list := make(Announcements, 0)

	var opts announcementsOptions
	for _, option := range options {
		option(&opts)
	}

	query := announcementsApplyFilters(
		r.db.SelectFrom("announcements"),
		opts.AnnouncementsFilters,
	)

	if opts.Pagination.IsZero() {
		return list, nil, handleErr(query.All(&list))
	}

	pageQuery := query.
		Paginate(opts.Pagination.Count).
		Cursor("-id")

	if err := paginate(pageQuery, opts.Pagination).All(&list); err != nil {
		return nil, nil, err
	}

	pi, err := r.pagesInfo(pageQuery, opts.Pagination, list)

	if err != nil {
		return nil, nil, err
	}

	return list, pi, nil
}

func (r *Repository) AnnouncementsOneByID(id int64) (*Announcement, error) {
	var o Announcement

	err := r.db.
		SelectFrom("announcements").
		Where("id", id).
		Limit(1).
		One(&o)

	if err != nil {
		return nil, handleErr(err)
	}

	return &o, nil
}
