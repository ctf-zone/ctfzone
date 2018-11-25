package controllers_test

import (
	"testing"
)

func TestGameInfo(t *testing.T) {
	setup(t)
	defer teardown(t)

	e := heDefault(t)
	e = heHost(e, "ctfzone.test")
	e = heCSRF(e)

	res := e.GET("/api/game").
		Expect().
		Status(200)

	checkJSONSchema(t, "GameInfo.json", res.Body().Raw())
}
