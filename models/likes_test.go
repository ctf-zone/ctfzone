package models_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	. "github.com/ctf-zone/ctfzone/models"
)

func Test_Likes_Insert_Success(t *testing.T) {
	setup(t)
	defer teardown(t)

	l := &Like{
		UserID:      1,
		ChallengeID: 2,
	}

	err := db.LikesInsert(l)
	assert.NoError(t, err)

	assert.WithinDuration(t, time.Now().UTC(), l.CreatedAt, 5*time.Second)
}

func Test_Likes_Insert_Duplicate(t *testing.T) {
	setup(t)
	defer teardown(t)

	l := &Like{
		UserID:      1,
		ChallengeID: 1,
	}

	err := db.LikesInsert(l)
	assert.Error(t, err)
}

func Test_Likes_Delete_Success(t *testing.T) {
	setup(t)
	defer teardown(t)

	err := db.LikesDelete(1, 1)
	assert.NoError(t, err)
}

func Test_Likes_Delete_NotExist(t *testing.T) {
	setup(t)
	defer teardown(t)

	err := db.LikesDelete(1337, 1337)
	assert.Error(t, err)
	assert.EqualError(t, err, sql.ErrNoRows.Error())
}

func Test_Likes_OneByID_Success(t *testing.T) {
	setup(t)
	defer teardown(t)

	s, err := db.LikesOneByID(1, 1)
	assert.NoError(t, err)
	assert.NotNil(t, s)
}

func Test_Likes_OneByID_NotExist(t *testing.T) {
	setup(t)
	defer teardown(t)

	s, err := db.LikesOneByID(1, 1337)
	assert.Error(t, err)
	assert.Nil(t, s)
	assert.EqualError(t, err, sql.ErrNoRows.Error())
}
