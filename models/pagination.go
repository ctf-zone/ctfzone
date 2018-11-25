package models

import (
	udb "upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

type Pagination struct {
	Count  uint  `schema:"count"`
	After  int64 `schema:"after"`
	Before int64 `schema:"before"`
}

var defaultPagination = Pagination{Count: 30}

func (p *Pagination) IsZero() bool {
	return p.Count == 0 && p.After == 0 && p.Before == 0
}

func (p *Pagination) IsForward() bool {
	return p.After != 0
}

func (p *Pagination) IsBackward() bool {
	return p.Before != 0
}

type PagesInfo struct {
	Next Pagination
	Prev Pagination
}

type pageable interface {
	First() int64
	Last() int64
	Len() int
}

func paginate(query sqlbuilder.Paginator, p Pagination) sqlbuilder.Paginator {

	if p.After != 0 {
		query = query.NextPage(p.After)
	} else if p.Before != 0 {
		query = query.PrevPage(p.Before)
	}

	return query
}

func (r *Repository) count(query sqlbuilder.Getter) (int, error) {
	var count int

	row, err := r.db.
		SelectFrom(udb.Raw("?", query)).
		As("page").
		Columns(udb.Raw("COUNT(*)")).
		QueryRow()

	if err != nil {
		return 0, err
	}

	if err := row.Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func (r *Repository) pagesInfo(query sqlbuilder.Paginator, p Pagination, models pageable) (*PagesInfo, error) {

	pi := &PagesInfo{}

	var first, last int64

	if models.Len() > 0 {
		first, last = models.First(), models.Last()
	} else if p.IsForward() {
		first, last = p.After, p.After
	} else {
		first, last = p.Before, p.Before
	}

	if count, err := r.count(query.NextPage(last)); err != nil {
		return nil, err
	} else if count > 0 {
		pi.Next.After = last
		pi.Next.Count = p.Count
	}

	if count, err := r.count(query.PrevPage(first)); err != nil {
		return nil, err
	} else if count > 0 {
		pi.Prev.Before = first
		pi.Prev.Count = p.Count
	}

	return pi, nil
}
