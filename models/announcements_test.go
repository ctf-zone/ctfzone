package models_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	. "github.com/ctf-zone/ctfzone/models"
)

func Test_Announcements_Insert_Success(t *testing.T) {
	setup(t)
	defer teardown(t)

	o := &Announcement{
		Title:       "Test title",
		Body:        "Test body",
		ChallengeID: nil,
	}

	err := db.AnnouncementsInsert(o)
	assert.NoError(t, err)
	assert.NotZero(t, o.ID)
	assert.WithinDuration(t, time.Now().UTC(), o.CreatedAt, 5*time.Second)
	assert.WithinDuration(t, time.Now().UTC(), o.UpdatedAt, 5*time.Second)
}

func Test_Announcements_Update_Success(t *testing.T) {
	setup(t)
	defer teardown(t)

	o1, err := db.AnnouncementsOneByID(1)
	assert.NoError(t, err)

	o1.Title = o1.Title + " [Updated]"
	updatedAt := o1.UpdatedAt

	err = db.AnnouncementsUpdate(o1)
	assert.NoError(t, err)

	o2, err := db.AnnouncementsOneByID(1)
	assert.NoError(t, err)
	assert.True(t, o2.UpdatedAt.After(updatedAt))
	assert.WithinDuration(t, time.Now(), o2.UpdatedAt, 5*time.Second)

	assert.Equal(t, o1, o2)
}

func Test_Announcements_Update_NotExist(t *testing.T) {
	setup(t)
	defer teardown(t)

	o := &Announcement{
		ID: 1337,
	}

	err := db.AnnouncementsUpdate(o)
	assert.Error(t, err)
	assert.EqualError(t, err, sql.ErrNoRows.Error())
}

func Test_Announcements_OneByID_Success(t *testing.T) {
	o, err := db.AnnouncementsOneByID(1)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), o.ID)
}

func Test_Announcements_OneByID_NotExist(t *testing.T) {
	o, err := db.AnnouncementsOneByID(1337)
	assert.Error(t, err)
	assert.EqualError(t, err, sql.ErrNoRows.Error())
	assert.Nil(t, o)
}

func Test_Announcements_Delete_Success(t *testing.T) {
	setup(t)
	defer teardown(t)

	err := db.AnnouncementsDelete(1)
	assert.NoError(t, err)

	o, err := db.AnnouncementsOneByID(1)
	assert.Error(t, err)
	assert.EqualError(t, err, sql.ErrNoRows.Error())
	assert.Nil(t, o)
}

func Test_Announcements_Delete_NotExist(t *testing.T) {
	setup(t)
	defer teardown(t)

	err := db.AnnouncementsDelete(1337)
	assert.Error(t, err)
	assert.EqualError(t, err, sql.ErrNoRows.Error())
}

func Test_Announcements_List_Basic(t *testing.T) {
	setup(t)
	defer teardown(t)

	list, _, err := db.AnnouncementsList(
		AnnouncementsListPagination(Pagination{Count: 10}),
	)
	assert.NoError(t, err)
	assert.Len(t, list, 3)
}
