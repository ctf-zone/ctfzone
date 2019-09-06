package models

import (
	"fmt"
	"strings"
)

type Pagination struct {
	Count  uint  `schema:"count"`
	After  int64 `schema:"after"`
	Before int64 `schema:"before"`
}

var defaultPagination = Pagination{Count: 2}

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

func paginate(query string, params map[string]interface{}, p Pagination) (string, map[string]interface{}) {

	if p.After != 0 {
		query, params = nextPage(query, params, p.After)
	} else if p.Before != 0 {
		query, params = prevPage(query, params, p.Before)
	}

	query += " ORDER BY id DESC"

	if p.Count != 0 {
		query += " LIMIT :count"
		params["count"] = p.Count
	}

	return query, params
}

func condPrefix(query string) string {
	var prefix string
	if strings.Contains(query, "WHERE") {
		prefix = " AND"
	} else {
		prefix = " WHERE"
	}
	return prefix
}

func nextPage(query string, params map[string]interface{}, after int64) (string, map[string]interface{}) {
	query += condPrefix(query) + " id < :after"
	params["after"] = after
	return query, params
}

func prevPage(query string, params map[string]interface{}, before int64) (string, map[string]interface{}) {
	query = fmt.Sprintf("SELECT * FROM (%s%s id > :before ORDER BY id ASC) AS p", query, condPrefix(query))
	params["before"] = before
	return query, params
}

func (r *Repository) count(query string, params map[string]interface{}) (int, error) {

	stmt, err := r.db.PrepareNamed(fmt.Sprintf("SELECT COUNT(*) FROM (%s) AS c", query))
	if err != nil {
		return 0, err
	}

	var count int

	if err := stmt.QueryRow(params).Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func (r *Repository) pagesInfo(query string, params map[string]interface{}, p Pagination, models pageable) (*PagesInfo, error) {
	pi := &PagesInfo{}

	var first, last int64

	if models.Len() > 0 {
		first, last = models.First(), models.Last()
	} else if p.IsForward() {
		first, last = p.After, p.After
	} else {
		first, last = p.Before, p.Before
	}

	if count, err := r.count(nextPage(query, params, last)); err != nil {
		return nil, err
	} else if count > 0 {
		pi.Next.After = last
		pi.Next.Count = p.Count
	}

	if count, err := r.count(prevPage(query, params, first)); err != nil {
		return nil, err
	} else if count > 0 {
		pi.Prev.Before = first
		pi.Prev.Count = p.Count
	}

	return pi, nil
}
