package models

import (
	"upper.io/db.v3/lib/sqlbuilder"
	"upper.io/db.v3/postgresql"
)

type Repository struct {
	db sqlbuilder.Database
}

func New(dsn string) (*Repository, error) {

	conn, err := postgresql.ParseURL(dsn)
	if err != nil {
		return nil, err
	}

	db, err := postgresql.Open(conn)
	if err != nil {
		return nil, err
	}

	return NewWithDB(db)
}

func NewWithDB(db sqlbuilder.Database) (*Repository, error) {
	db.SetPreparedStatementCache(true)

	// TODO: set logger to logrus

	return &Repository{db}, nil
}

func (r *Repository) Close() error {
	return r.db.Close()
}

func (r *Repository) EnableLogging() {
	r.db.SetLogging(true)
}
