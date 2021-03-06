package models

import (
	"time"
)

type Solution struct {
	UserID      int64     `db:"user_id"      json:"userId"`
	ChallengeID int64     `db:"challenge_id" json:"challengeId"`
	CreatedAt   time.Time `db:"created_at"   json:"createdAt"`
}

func (r *Repository) SolutionsInsert(o *Solution) error {
	o.CreatedAt = now()

	stmt, err := r.db.PrepareNamed(
		"INSERT INTO solutions (user_id, challenge_id, created_at) " +
			"VALUES(:user_id, :challenge_id, :created_at) " +
			"RETURNING created_at")

	if err != nil {
		return err
	}

	return stmt.QueryRowx(o).Scan(&o.CreatedAt)
}

func (r *Repository) SolutionsOneByID(userID, challengeID int64) (*Solution, error) {
	var o Solution

	err := r.db.Get(&o, "SELECT * FROM solutions WHERE user_id = $1 and challenge_id = $2", userID, challengeID)

	if err != nil {
		return nil, err
	}

	return &o, nil
}
