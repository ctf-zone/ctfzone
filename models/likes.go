package models

import (
	"database/sql"
	"time"
)

type LikesRepository interface {
	Insert(*Like) error
	Delete(int64, int64) error

	OneByID(int64, int64) (*Like, error)
}

type Like struct {
	UserID      int64     `db:"user_id"      json:"userId"`
	ChallengeID int64     `db:"challenge_id" json:"challengeId"`
	CreatedAt   time.Time `db:"created_at"   json:"createdAt"`
}

func (r *Repository) LikesInsert(o *Like) error {

	o.CreatedAt = now()

	// TODO: id
	_, err := r.db.
		InsertInto("likes").
		Values(o).
		Exec()

	return err
}

func (r *Repository) LikesDelete(userID, challengeID int64) error {
	res, err := r.db.
		DeleteFrom("likes").
		Where("user_id", userID).
		And("challenge_id", challengeID).
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

func (r *Repository) LikesOneByID(userID, challengeID int64) (*Like, error) {
	var o Like

	err := r.db.
		SelectFrom("likes").
		Where("user_id", userID).
		And("challenge_id", challengeID).
		One(&o)

	if err != nil {
		return nil, handleErr(err)
	}

	return &o, nil
}
