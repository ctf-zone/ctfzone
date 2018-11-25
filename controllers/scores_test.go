package controllers_test

import (
	"testing"
)

func TestScoresList(t *testing.T) {
	setup(t)
	defer teardown(t)

	e := heDefault(t)
	e = heHost(e, "ctfzone.test")

	res := e.GET("/api/scores").
		Expect().
		Status(200)

	checkJSONSchema(t, "Scores.json", res.Body().Raw())
}

func TestScoresCtftimeList(t *testing.T) {
	setup(t)
	defer teardown(t)

	e := heDefault(t)
	e = heHost(e, "ctfzone.test")

	res := e.GET("/api/scores/ctftime").
		Expect().
		Status(200)

	checkJSONSchema(t, "ScoresCTFTime.json", res.Body().Raw())
}
