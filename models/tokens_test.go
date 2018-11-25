package models_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	. "github.com/ctf-zone/ctfzone/models"
)

func Test_Tokens_Insert_Success(t *testing.T) {
	setup(t)
	defer teardown(t)

	o, err := db.TokensNew(1, TokenTypeActivate, 24*time.Hour)
	assert.NoError(t, err)

	assert.WithinDuration(t, time.Now().Add(24*time.Hour).UTC(), o.ExpiresAt, 5*time.Second)

	err = db.TokensInsert(o)
	assert.NoError(t, err)
	assert.NotZero(t, o.ID)
	assert.WithinDuration(t, time.Now().UTC(), o.CreatedAt, 5*time.Second)
}

func Test_Tokens_Insert_UserNotExist(t *testing.T) {
	setup(t)
	defer teardown(t)

	o, err := db.TokensNew(1337, TokenTypeActivate, 24*time.Hour)
	assert.NoError(t, err)
	assert.WithinDuration(t, time.Now().Add(24*time.Hour).UTC(), o.ExpiresAt, 5*time.Second)

	err = db.TokensInsert(o)
	assert.Error(t, err)
}

func Test_Tokens_OneByID_Success(t *testing.T) {
	setup(t)
	defer teardown(t)

	o, err := db.TokensOneByID(2)
	assert.NoError(t, err)
	assert.NotNil(t, o)
	assert.Equal(t, int64(2), o.ID)
}

func Test_Tokens_OneByID_NotExist(t *testing.T) {
	setup(t)
	defer teardown(t)

	o, err := db.TokensOneByID(1337)
	assert.Error(t, err)
	assert.EqualError(t, err, sql.ErrNoRows.Error())
	assert.Nil(t, o)
}

func Test_Tokens_OneByTokenAndType_Success(t *testing.T) {
	setup(t)
	defer teardown(t)

	o, err := db.TokensOneByTokenAndType(
		"1534e7c3fc6b3dbb411a5170b5fa94cd95f324f1ca072853ab2cb34c1378c061",
		TokenTypeActivate,
	)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), o.ID)
}

func Test_Tokens_OneByTokenAndType_NotExist(t *testing.T) {
	setup(t)
	defer teardown(t)

	o, err := db.TokensOneByTokenAndType(
		"2534e7c3fc6b3dbb411a5170b5fa94cd95f324f1ca072853ab2cb34c1378c061",
		TokenTypeActivate,
	)
	assert.Error(t, err)
	assert.Nil(t, o)
}

func Test_Tokens_OneByUserAndType_Success(t *testing.T) {
	setup(t)
	defer teardown(t)

	o, err := db.TokensOneByUserAndType(3, TokenTypeReset)
	assert.NoError(t, err)
	assert.Equal(t, int64(3), o.ID)
}

func Test_Tokens_OneByUserAndType_NotExist(t *testing.T) {
	setup(t)
	defer teardown(t)

	o, err := db.TokensOneByUserAndType(1, TokenTypeReset)

	assert.Error(t, err)
	assert.EqualError(t, err, sql.ErrNoRows.Error())
	assert.Nil(t, o)
}

func Test_Tokens_Delete_Success(t *testing.T) {
	setup(t)
	defer teardown(t)

	err := db.TokensDelete(1)
	assert.NoError(t, err)

	o, err := db.TokensOneByID(1)
	assert.Error(t, err)
	assert.EqualError(t, err, sql.ErrNoRows.Error())
	assert.Nil(t, o)
}

func Test_Tokens_Delete_NotExist(t *testing.T) {
	setup(t)
	defer teardown(t)

	err := db.TokensDelete(1337)
	assert.Error(t, err)
	assert.EqualError(t, err, sql.ErrNoRows.Error())
}
