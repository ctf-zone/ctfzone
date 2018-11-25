package controllers_test

import (
	"testing"

	. "github.com/ctf-zone/ctfzone/controllers"
)

func TestChallengesList(t *testing.T) {
	setup(t)
	defer teardown(t)

	e := heDefault(t)
	e = heHost(e, "ctfzone.test")
	e = heCSRF(e)
	e = heAuth(e, "lcbc@mail.com", "lcbc")

	res := e.GET("/api/challenges").
		Expect().
		Status(200)

	checkJSONSchema(t, "Challenges.json", res.Body().Raw())

	res.JSON().Array().Length().Equal(3)
}

func TestChallengesGet_InvalidID(t *testing.T) {
	setup(t)
	defer teardown(t)

	e := heDefault(t)
	e = heHost(e, "ctfzone.test")
	e = heCSRF(e)
	e = heAuth(e, "lcbc@mail.com", "lcbc")

	res := e.GET("/api/challenges/<invalid>").
		Expect().
		Status(400)

	res.JSON().Path("$.error").String().Equal(ErrInvalidID.Msg)
}

func TestChallengesGet_NotExistingChallenge(t *testing.T) {
	setup(t)
	defer teardown(t)

	e := heDefault(t)
	e = heHost(e, "ctfzone.test")
	e = heCSRF(e)
	e = heAuth(e, "lcbc@mail.com", "lcbc")

	res := e.GET("/api/challenges/1337").
		Expect().
		Status(404)

	res.JSON().Path("$.error").String().Equal(ErrChallengeNotFound.Msg)
}

func TestChallengesGet_Success(t *testing.T) {
	setup(t)
	defer teardown(t)

	e := heDefault(t)
	e = heHost(e, "ctfzone.test")
	e = heCSRF(e)
	e = heAuth(e, "lcbc@mail.com", "lcbc")

	res := e.GET("/api/challenges/1").
		Expect().
		Status(200)

	checkJSONSchema(t, "ChallengeE.json", res.Body().Raw())
}
