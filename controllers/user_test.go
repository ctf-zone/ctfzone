package controllers_test

import (
	"testing"

	. "github.com/ctf-zone/ctfzone/controllers"
)

func TestUserSolutionsCreate_InvalidChallengeID(t *testing.T) {
	setup(t)
	defer teardown(t)

	e := heDefault(t)
	e = heHost(e, "ctfzone.test")
	e = heCSRF(e)
	e = heAuth(e, "lcbc@mail.com", "lcbc")

	res := e.POST("/api/user/solutions/wrong").
		WithJSON(map[string]interface{}{
			"flag": "ctfzone{reverse-200}",
		}).
		Expect().
		Status(400)

	res.JSON().Path("$.error").String().Equal(ErrInvalidID.Msg)
}

func TestUserSolutionsCreate_NotExistingChallenge(t *testing.T) {
	setup(t)
	defer teardown(t)

	e := heDefault(t)
	e = heHost(e, "ctfzone.test")
	e = heCSRF(e)
	e = heAuth(e, "lcbc@mail.com", "lcbc")

	res := e.POST("/api/user/solutions/1337").
		WithJSON(map[string]interface{}{
			"flag": "ctfzone{reverse-200}",
		}).
		Expect().
		Status(404)

	res.JSON().Path("$.error").String().Equal(ErrChallengeNotFound.Msg)
}

func TestUserSolutionsCreate_WrongFlag(t *testing.T) {
	setup(t)
	defer teardown(t)

	e := heDefault(t)
	e = heHost(e, "ctfzone.test")
	e = heCSRF(e)
	e = heAuth(e, "lcbc@mail.com", "lcbc")

	res := e.POST("/api/user/solutions/2").
		WithJSON(map[string]interface{}{
			"flag": "ctfzone{wrong}",
		}).
		Expect().
		Status(418)

	res.JSON().Path("$.error").String().Equal(ErrInvalidFlag.Msg)
}

func TestUserSolutionsCreate_Success(t *testing.T) {
	setup(t)
	defer teardown(t)

	e := heDefault(t)
	e = heHost(e, "ctfzone.test")
	e = heCSRF(e)
	e = heAuth(e, "lcbc@mail.com", "lcbc")

	e.POST("/api/user/solutions/2").
		WithJSON(map[string]interface{}{
			"flag": "ctfzone{reverse-200}",
		}).
		Expect().
		Status(200).NoContent()

	// Check solution exists
	// TODO
}

func TestUserLikesCreate_NotExistingChallenge(t *testing.T) {
	setup(t)
	defer teardown(t)

	e := heDefault(t)
	e = heHost(e, "ctfzone.test")
	e = heCSRF(e)
	e = heAuth(e, "lcbc@mail.com", "lcbc")

	res := e.POST("/api/user/likes/1337").
		WithHeader("Content-Type", "application/json").
		WithText("{}").
		Expect().
		Status(404)

	res.JSON().Path("$.error").String().Equal(ErrChallengeNotFound.Msg)
}

func TestUserLikesCreate_InvalidChallengeID(t *testing.T) {
	setup(t)
	defer teardown(t)

	e := heDefault(t)
	e = heHost(e, "ctfzone.test")
	e = heCSRF(e)
	e = heAuth(e, "lcbc@mail.com", "lcbc")

	res := e.POST("/api/user/likes/wrong").
		WithHeader("Content-Type", "application/json").
		WithText("{}").
		Expect().
		Status(400)

	res.JSON().Path("$.error").String().Equal(ErrInvalidID.Msg)
}

func TestUserLikesCreate_AlreadyLiked(t *testing.T) {
	setup(t)
	defer teardown(t)

	e := heDefault(t)
	e = heHost(e, "ctfzone.test")
	e = heCSRF(e)
	e = heAuth(e, "lcbc@mail.com", "lcbc")

	res := e.POST("/api/user/likes/1").
		WithHeader("Content-Type", "application/json").
		WithText("{}").
		Expect().
		Status(409)

	res.JSON().Path("$.error").String().Contains(ErrDuplicate.Msg)
}

func TestUserLikesCreate_Success(t *testing.T) {
	setup(t)
	defer teardown(t)

	e := heDefault(t)
	e = heHost(e, "ctfzone.test")
	e = heCSRF(e)
	e = heAuth(e, "lcbc@mail.com", "lcbc")

	e.POST("/api/user/likes/2").
		WithHeader("Content-Type", "application/json").
		WithText("{}").
		Expect().
		Status(201).NoContent()

	// TODO: check exist
}

func TestUserLikesDelete_InvalidChallengeID(t *testing.T) {
	setup(t)
	defer teardown(t)

	e := heDefault(t)
	e = heHost(e, "ctfzone.test")
	e = heCSRF(e)
	e = heAuth(e, "lcbc@mail.com", "lcbc")

	res := e.DELETE("/api/user/likes/wrong").
		Expect().
		Status(400)

	res.JSON().Path("$.error").String().Equal(ErrInvalidID.Msg)
}

func TestUserLikesDelete_NotExistingChallenge(t *testing.T) {
	setup(t)
	defer teardown(t)

	e := heDefault(t)
	e = heHost(e, "ctfzone.test")
	e = heCSRF(e)
	e = heAuth(e, "lcbc@mail.com", "lcbc")

	res := e.DELETE("/api/user/likes/1337").
		Expect().
		Status(404)

	res.JSON().Path("$.error").String().Equal(ErrChallengeNotFound.Msg)
}

func TestUserLikesDelete_Success(t *testing.T) {
	setup(t)
	defer teardown(t)

	e := heDefault(t)
	e = heHost(e, "ctfzone.test")
	e = heCSRF(e)
	e = heAuth(e, "lcbc@mail.com", "lcbc")

	e.DELETE("/api/user/likes/1").
		Expect().
		Status(200).NoContent()

	// TODO: check exist
}
