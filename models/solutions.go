package models

import (
	"time"
)

type SolutionsRepository interface {
	Insert(*Solution) error

	OneByID(int64, int64) (*Solution, error)
}

type Solution struct {
	UserID      int64     `db:"user_id"      json:"userId"`
	ChallengeID int64     `db:"challenge_id" json:"challengeId"`
	CreatedAt   time.Time `db:"created_at"   json:"createdAt"`
}

func (r *Repository) SolutionsInsert(o *Solution) error {
	o.CreatedAt = now()

	// TODO: id
	_, err := r.db.
		InsertInto("solutions").
		Values(o).
		Exec()

	return err
}

func (r *Repository) SolutionsOneByID(userID, challengeID int64) (*Solution, error) {
	var o Solution

	err := r.db.
		SelectFrom("solutions").
		Where("user_id", userID).
		And("challenge_id", challengeID).
		One(&o)

	if err != nil {
		return nil, handleErr(err)
	}

	return &o, nil
}
