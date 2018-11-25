package controllers_test

import (
	"testing"
)

func TestAnnouncementsList_Success(t *testing.T) {
	setup(t)
	defer teardown(t)

	e := heDefault(t)
	e = heHost(e, "ctfzone.test")
	e = heCSRF(e)

	res := e.GET("/api/announcements").
		Expect().
		Status(200)

	checkJSONSchema(t, "Announcements.json", res.Body().Raw())
}
