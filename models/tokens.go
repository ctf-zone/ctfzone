package models

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"time"
)

type TokenType string

const (
	TokenTypeActivate = TokenType("activate")
	TokenTypeReset    = TokenType("reset")
)

type Token struct {
	ID        int64     `db:"id,omitempty" json:"id"`
	UserID    int64     `db:"user_id"      json:"userId"`
	Token     string    `db:"token"        json:"token"`
	Type      TokenType `db:"type"         json:"type"`
	ExpiresAt time.Time `db:"expires_at"   json:"expiresAt"`
	CreatedAt time.Time `db:"created_at"   json:"createdAt"`
}

const tokenLength = 32

func (r *Repository) TokensNew(userID int64, tp TokenType, lifetime time.Duration) (*Token, error) {
	token := make([]byte, tokenLength)

	if n, err := io.ReadFull(rand.Reader, token); err != nil || n != tokenLength {
		return nil, err
	}

	t := &Token{
		UserID:    userID,
		Type:      tp,
		Token:     hex.EncodeToString(token),
		ExpiresAt: time.Now().Add(lifetime).UTC(),
	}

	return t, nil
}

func (r *Repository) TokensInsert(o *Token) error {
	o.CreatedAt = now()

	stmt, err := r.db.PrepareNamed(
		"INSERT INTO tokens (user_id, token, type, expires_at, created_at) " +
			"VALUES(:user_id, :token, :type, :expires_at, :created_at) " +
			"RETURNING id")

	if err != nil {
		return err
	}

	return stmt.QueryRowx(o).Scan(&o.ID)
}

func (r *Repository) TokensOneByID(id int64) (*Token, error) {
	var o Token

	err := r.db.Get(&o, "SELECT * FROM tokens WHERE id = $1", id)

	if err != nil {
		return nil, err
	}

	return &o, nil
}

func (r *Repository) TokensOneByTokenAndType(token string, tp TokenType) (*Token, error) {
	var o Token

	err := r.db.Get(&o, "SELECT * FROM tokens WHERE token = $1 and type = $2", token, tp)

	if err != nil {
		return nil, err
	}

	return &o, nil
}

func (r *Repository) TokensOneByUserAndType(userID int64, tp TokenType) (*Token, error) {
	var o Token

	err := r.db.Get(&o, "SELECT * FROM tokens WHERE user_id = $1 and type = $2", userID, tp)

	if err != nil {
		return nil, err
	}

	return &o, nil
}

func (r *Repository) TokensDelete(id int64) error {
	return r.db.QueryRow("DELETE FROM tokens WHERE id = $1 RETURNING id", id).Scan(&id)
}
