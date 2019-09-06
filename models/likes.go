package models

import (
	"time"
)

type Like struct {
	UserID      int64     `db:"user_id"      json:"userId"`
	ChallengeID int64     `db:"challenge_id" json:"challengeId"`
	CreatedAt   time.Time `db:"created_at"   json:"createdAt"`
}

func (r *Repository) LikesInsert(o *Like) error {
	o.CreatedAt = now()

	stmt, err := r.db.PrepareNamed(
		"INSERT INTO likes (user_id, challenge_id, created_at) " +
			"VALUES(:user_id, :challenge_id, :created_at) " +
			"RETURNING created_at")

	if err != nil {
		return err
	}

	return stmt.QueryRowx(o).Scan(&o.CreatedAt)
}

func (r *Repository) LikesDelete(userID, challengeID int64) error {
	return r.db.QueryRow("DELETE FROM likes WHERE user_id = $1 AND challenge_id = $2 RETURNING user_id",
		userID, challengeID).Scan(&userID)
}

func (r *Repository) LikesOneByID(userID, challengeID int64) (*Like, error) {
	var o Like

	err := r.db.Get(&o, "SELECT * FROM likes WHERE user_id = $1 and challenge_id = $2", userID, challengeID)

	if err != nil {
		return nil, err
	}

	return &o, nil
}
