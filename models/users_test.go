package models_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	. "github.com/ctf-zone/ctfzone/models"
)

func Test_Users_Insert_Success(t *testing.T) {
	setup(t)
	defer teardown(t)

	o := &User{
		Name:     "team",
		Email:    "team@mail.com",
		Password: "12345678",
		Extra: map[string]interface{}{
			"country": "RU",
		},
		IsActivated: false,
	}

	err := db.UsersInsert(o)
	assert.NoError(t, err)
	assert.NotZero(t, o.ID)
	assert.WithinDuration(t, time.Now().UTC(), o.CreatedAt, 5*time.Second)
	assert.WithinDuration(t, time.Now().UTC(), o.UpdatedAt, 5*time.Second)
	assert.NoError(t,
		bcrypt.CompareHashAndPassword(
			[]byte(o.PasswordHash),
			[]byte("12345678"),
		),
	)
}

func Test_Users_Insert_Duplicate(t *testing.T) {
	setup(t)
	defer teardown(t)

	u := &User{
		Name:     "ppp",
		Email:    "ppp@mail.com",
		Password: "12345678",
		Extra: map[string]interface{}{
			"country": "US",
		},
		IsActivated: false,
	}

	err := db.UsersInsert(u)
	assert.Error(t, err)
}

func Test_Users_Update_Success(t *testing.T) {
	setup(t)
	defer teardown(t)

	o1, err := db.UsersOneByID(1)
	assert.NoError(t, err)

	o1.IsActivated = true
	updatedAt := o1.UpdatedAt

	err = db.UsersUpdate(o1)
	assert.NoError(t, err)

	o2, err := db.UsersOneByID(1)
	assert.NoError(t, err)

	assert.True(t, o2.UpdatedAt.After(updatedAt))
	assert.WithinDuration(t, time.Now(), o2.UpdatedAt, 5*time.Second)

	assert.Equal(t, o1, o2)
}

func Test_Users_Update_ChangePassword(t *testing.T) {
	setup(t)
	defer teardown(t)

	o, err := db.UsersOneByEmail("lcbc@mail.com")
	assert.NoError(t, err)

	o.Password = "new-password"

	err = db.UsersUpdate(o)
	assert.NoError(t, err)

	_, err = db.UsersLogin("lcbc@mail.com", "new-password")
	assert.NoError(t, err)
}

func Test_Users_Update_EmptyPassword(t *testing.T) {
	setup(t)
	defer teardown(t)

	o, err := db.UsersOneByEmail("lcbc@mail.com")
	assert.NoError(t, err)

	o.Password = ""

	err = db.UsersUpdate(o)
	assert.NoError(t, err)

	_, err = db.UsersLogin("lcbc@mail.com", "lcbc")
	assert.NoError(t, err)
}

func Test_Users_Update_NotExist(t *testing.T) {
	setup(t)
	defer teardown(t)

	o := &User{
		ID:   1337,
		Name: "test",
	}

	err := db.UsersUpdate(o)
	assert.Error(t, err)
	assert.EqualError(t, err, sql.ErrNoRows.Error())
}

func Test_Users_OneByID_Success(t *testing.T) {
	setup(t)
	defer teardown(t)

	o, err := db.UsersOneByID(1)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), o.ID)
	assert.Equal(t, "lcbc@mail.com", o.Email)
}

func Test_Users_OneByID_NotExist(t *testing.T) {
	setup(t)
	defer teardown(t)

	o, err := db.UsersOneByID(1337)
	assert.Error(t, err)
	assert.EqualError(t, err, sql.ErrNoRows.Error())
	assert.Nil(t, o)
}

func Test_Users_OneByEmail_Success(t *testing.T) {
	setup(t)
	defer teardown(t)

	o, err := db.UsersOneByEmail("lcbc@mail.com")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), o.ID)
	assert.Equal(t, "lcbc@mail.com", o.Email)
}

func Test_Users_OneByEmail_NotExist(t *testing.T) {
	setup(t)
	defer teardown(t)

	o, err := db.UsersOneByEmail("not-exist@mail.com")
	assert.Error(t, err)
	assert.EqualError(t, err, sql.ErrNoRows.Error())
	assert.Nil(t, o)
}

func Test_Users_OneByName_Success(t *testing.T) {
	setup(t)
	defer teardown(t)

	o, err := db.UsersOneByName("PPP")
	assert.NoError(t, err)

	assert.Equal(t, int64(2), o.ID)
	assert.Equal(t, "PPP", o.Name)
}

func Test_Users_OneByName_NotExist(t *testing.T) {
	setup(t)
	defer teardown(t)

	_, err := db.UsersOneByName("NotExist")
	assert.Error(t, err)
	assert.EqualError(t, err, sql.ErrNoRows.Error())
}

func Test_Users_Login_Success(t *testing.T) {
	setup(t)
	defer teardown(t)

	o, err := db.UsersLogin("lcbc@mail.com", "lcbc")
	assert.NoError(t, err)
	assert.NotZero(t, o.ID)
}

func Test_Users_Login_NotExist(t *testing.T) {
	setup(t)
	defer teardown(t)

	_, err := db.UsersLogin("not-exist@mail.com", "lcbc")
	assert.Error(t, err)
}

func Test_UsersLogin_WrongPassword(t *testing.T) {
	setup(t)
	defer teardown(t)

	_, err := db.UsersLogin("lcbc@mail.com", "1234")
	assert.Error(t, err)
}

func Test_Users_List_FirstPage(t *testing.T) {
	setup(t)
	defer teardown(t)

	users, _, err := db.UsersList(
		UsersPagination(Pagination{Count: 3}),
	)
	require.NoError(t, err)

	require.Len(t, users, 3)
	assert.Equal(t, int64(5), users[0].ID)
	assert.Equal(t, int64(3), users[2].ID)

	// assert.True(t, pi.HasNext)
	// assert.False(t, pi.HasPrev)
}

func Test_Users_List_NextPage(t *testing.T) {
	setup(t)
	defer teardown(t)

	users, _, err := db.UsersList(
		UsersPagination(Pagination{
			Count: 3,
			After: 3,
		}),
	)
	require.NoError(t, err)

	require.Len(t, users, 2)
	assert.Equal(t, int64(2), users[0].ID)
	assert.Equal(t, int64(1), users[1].ID)

	// assert.False(t, pi.HasNext)
	// assert.True(t, pi.HasPrev)
}

func Test_Users_List_PrevPage(t *testing.T) {
	setup(t)
	defer teardown(t)

	users, _, err := db.UsersList(
		UsersPagination(Pagination{
			Count:  3,
			Before: 3,
		}),
	)
	require.NoError(t, err)

	require.Len(t, users, 2)
	assert.Equal(t, int64(5), users[0].ID)
	assert.Equal(t, int64(4), users[1].ID)

	// assert.True(t, pi.HasNext)
	// assert.False(t, pi.HasPrev)
}

func Test_Users_Delete_Success(t *testing.T) {
	setup(t)
	defer teardown(t)

	err := db.UsersDelete(1)
	assert.NoError(t, err)
}

func Test_Users_Delete_NotExist(t *testing.T) {
	setup(t)
	defer teardown(t)

	err := db.UsersDelete(1337)
	assert.EqualError(t, err, sql.ErrNoRows.Error())
}
