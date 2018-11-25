package models

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"io"
	"time"
)

type TokensRepository interface {
	Insert(*Token) error
	Delete(int64) error

	OneByID(int64) (*Token, error)
	OneByTokenAndType(string, TokenType) (*Token, error)
	OneByUserAndType(int64, TokenType) (*Token, error)

	New(int64, TokenType, time.Duration) (*Token, error)
}

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

	row, err := r.db.
		InsertInto("tokens").
		Values(o).
		Returning("id").
		QueryRow()

	if err != nil {
		return err
	}

	return row.Scan(&o.ID)
}

func (r *Repository) TokensOneByID(id int64) (*Token, error) {
	var o Token

	err := r.db.
		SelectFrom("tokens").
		Where("id", id).
		Limit(1).
		One(&o)

	if err != nil {
		return nil, handleErr(err)
	}

	return &o, nil
}

func (r *Repository) TokensOneByTokenAndType(token string, tp TokenType) (*Token, error) {
	var o Token

	err := r.db.
		SelectFrom("tokens").
		Where("token", token).
		And("type", tp).
		Limit(1).
		One(&o)

	if err != nil {
		return nil, handleErr(err)
	}

	return &o, nil
}

func (r *Repository) TokensOneByUserAndType(userID int64, tp TokenType) (*Token, error) {
	var o Token

	err := r.db.
		SelectFrom("tokens").
		Where("user_id", userID).
		And("type", tp).
		Limit(1).
		One(&o)

	if err != nil {
		return nil, handleErr(err)
	}

	return &o, nil
}

func (r *Repository) TokensDelete(id int64) error {
	res, err := r.db.
		DeleteFrom("tokens").
		Where("id", id).
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
