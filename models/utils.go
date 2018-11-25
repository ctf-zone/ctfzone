package models

import (
	"database/sql"
	"time"

	udb "upper.io/db.v3"
)

func now() time.Time {
	// PostgreSQL currently can't store nanoseconds
	// https://github.com/lib/pq/issues/227
	return time.Now().UTC().Truncate(time.Microsecond)
}

func handleErr(err error) error {
	if err == udb.ErrNoMoreRows {
		return sql.ErrNoRows
	}

	return err
}
