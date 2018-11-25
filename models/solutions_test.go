package models_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	. "github.com/ctf-zone/ctfzone/models"
)

func Test_Solutions_Insert_Success(t *testing.T) {
	setup(t)
	defer teardown(t)

	o := &Solution{
		UserID:      5,
		ChallengeID: 2,
	}

	err := db.SolutionsInsert(o)
	assert.NoError(t, err)
	assert.WithinDuration(t, time.Now().UTC(), o.CreatedAt, 5*time.Second)
}

func Test_Solutions_Insert_Duplicate(t *testing.T) {
	setup(t)
	defer teardown(t)

	o := &Solution{
		UserID:      1,
		ChallengeID: 1,
	}

	err := db.SolutionsInsert(o)
	assert.Error(t, err)
}

func Test_Solutions_OneByID_Success(t *testing.T) {
	setup(t)
	defer teardown(t)

	o, err := db.SolutionsOneByID(1, 1)
	assert.NoError(t, err)
	assert.NotNil(t, o)
}

func Test_Solutions_OneByID_NotExist(t *testing.T) {
	setup(t)
	defer teardown(t)

	o, err := db.SolutionsOneByID(1, 1337)
	assert.Error(t, err)
	assert.Nil(t, o)
	assert.EqualError(t, err, sql.ErrNoRows.Error())
}
