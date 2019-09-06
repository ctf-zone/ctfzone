package models_test

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/ctf-zone/ctfzone/config"
	"github.com/ctf-zone/ctfzone/internal/crypto"
	. "github.com/ctf-zone/ctfzone/models"
)

func Test_Challenges_Insert_Success(t *testing.T) {
	setup(t)
	defer teardown(t)

	o := &Challenge{
		Title:       "Test challenge",
		Categories:  []string{"web"},
		Description: "Some description",
		Difficulty:  DifficultyMedium,
		Flag:        "ctfzone{test}",
		Points:      200,
		IsLocked:    false,
	}

	err := db.ChallengesInsert(o)
	assert.NoError(t, err)

	assert.NotZero(t, o.ID)
	assert.WithinDuration(t, time.Now().UTC(), o.CreatedAt, 5*time.Second)
	assert.WithinDuration(t, time.Now().UTC(), o.UpdatedAt, 5*time.Second)
	assert.Equal(t,
		o.FlagHash,
		fmt.Sprintf("%x", sha256.Sum256([]byte("ctfzone{test}"))),
	)
}

func Test_Challenges_Update_Success(t *testing.T) {
	setup(t)
	defer teardown(t)

	o1, err := db.ChallengesOneByID(1)
	assert.NoError(t, err)

	o1.Title = o1.Title + " [Updated]"
	updatedAt := o1.UpdatedAt

	err = db.ChallengesUpdate(o1)
	assert.NoError(t, err)

	o2, err := db.ChallengesOneByID(1)
	assert.NoError(t, err)
	assert.Equal(t, o1, o2)

	assert.True(t, o2.UpdatedAt.After(updatedAt))
	assert.WithinDuration(t, time.Now(), o2.UpdatedAt, 5*time.Second)
}

func Test_Challenges_Update_NotExist(t *testing.T) {
	setup(t)
	defer teardown(t)

	o := &Challenge{
		ID:         1337,
		Title:      "Test",
		Difficulty: DifficultyMedium,
	}

	err := db.ChallengesUpdate(o)
	assert.Error(t, err)
	assert.EqualError(t, err, sql.ErrNoRows.Error())
}

func Test_Challenges_Update_ChangeFlag(t *testing.T) {
	setup(t)
	defer teardown(t)

	o, err := db.ChallengesOneByID(1)
	assert.NoError(t, err)

	err = crypto.CheckFlag(o.FlagHash, "ctfzone{web-100}")
	assert.NoError(t, err)

	o.Flag = "ctfzone{updated}"

	err = db.ChallengesUpdate(o)
	assert.NoError(t, err)

	err = crypto.CheckFlag(o.FlagHash, "ctfzone{updated}")
	assert.NoError(t, err)
}

func Test_Challenges_Update_EmptyFlag(t *testing.T) {
	setup(t)
	defer teardown(t)

	o, err := db.ChallengesOneByID(1)
	assert.NoError(t, err)

	err = crypto.CheckFlag(o.FlagHash, "ctfzone{web-100}")
	assert.NoError(t, err)

	o.Flag = ""

	err = db.ChallengesUpdate(o)
	assert.NoError(t, err)

	err = crypto.CheckFlag(o.FlagHash, "ctfzone{web-100}")
	assert.NoError(t, err)
}

func Test_Challenges_OneByID_Success(t *testing.T) {
	o, err := db.ChallengesOneByID(1)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), o.ID)
	assert.Equal(t, "Web", o.Title)
}

func Test_Challenges_OneByID_NotExist(t *testing.T) {
	o, err := db.ChallengesOneByID(1337)
	assert.Error(t, err)
	assert.EqualError(t, err, sql.ErrNoRows.Error())
	assert.Nil(t, o)
}

func Test_Challenges_CheckFlag_Success(t *testing.T) {
	setup(t)
	defer teardown(t)

	o, err := db.ChallengesOneByID(1)
	assert.NoError(t, err)

	err = crypto.CheckFlag(o.FlagHash, "ctfzone{web-100}")
	assert.NoError(t, err)
}

func Test_Challenges_CheckFlag_WrongFlag(t *testing.T) {
	setup(t)
	defer teardown(t)

	o, err := db.ChallengesOneByID(1)
	assert.NoError(t, err)

	err = crypto.CheckFlag(o.Flag, "ctfzone{wrong}")
	assert.Error(t, err)
}

func Test_Challenges_ListExt_Basic(t *testing.T) {
	setup(t)
	defer teardown(t)

	list, _, err := db.ChallengesListE(
		&config.Scoring{},
		ChallengesPagination(Pagination{Count: 2}),
		ChallengesIncludeMeta(),
		ChallengesIncludeUser(1),
	)
	assert.NoError(t, err)
	assert.NotEmpty(t, list)
}
