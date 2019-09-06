package models

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Repository struct {
	db *sqlx.DB
}

func New(dsn string) (*Repository, error) {

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}

	return &Repository{db}, nil
}

func (r *Repository) Close() error {
	return r.db.Close()
}

func (r *Repository) EnableLogging() {
	// TODO
}
