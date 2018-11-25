package controllers_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"

	. "github.com/ctf-zone/ctfzone/controllers"
)

func TestAuthRegister_Duplicate(t *testing.T) {
	setup(t)
	defer teardown(t)

	e := heDefault(t)
	e = heHost(e, "ctfzone.test")
	e = heCSRF(e)

	res := e.POST("/api/auth/register").
		WithJSON(map[string]interface{}{
			"name":     "PPP",
			"email":    "ppp@mail.com",
			"password": "password",
			"extra":    json.RawMessage(`{"country": "RU"}`),
		}).
		Expect().
		Status(409)

	res.JSON().Path("$.error").String().Equal(ErrDuplicate.Msg)
	res.JSON().Path("$.errors").Object().Keys().ContainsOnly("name", "email")
}

func TestAuthRegister_Success(t *testing.T) {
	setup(t)
	defer teardown(t)

	e := heDefault(t)
	e = heHost(e, "ctfzone.test")
	e = heCSRF(e)

	done := make(chan struct{})

	m.On("Send",
		"activate",
		"team@mail.com",
		mock.Anything).
		Return(nil).
		Run(func(mock.Arguments) { close(done) })

	e.POST("/api/auth/register").
		WithJSON(map[string]interface{}{
			"name":     "team",
			"email":    "team@mail.com",
			"password": "password",
			"extra":    json.RawMessage(`{"country": "RU"}`),
		}).
		Expect().
		Status(201).
		NoContent()

	select {
	case <-done:
	case <-time.After(10 * time.Millisecond):
		t.Log("Timeout")
	}

	m.AssertExpectations(t)
}

func TestAuthLogin_InvalidCreds(t *testing.T) {
	setup(t)
	defer teardown(t)

	e := heDefault(t)
	e = heHost(e, "ctfzone.test")
	e = heCSRF(e)

	res := e.POST("/api/auth/login").
		WithJSON(map[string]interface{}{
			"email":    "team@mail.com",
			"password": "12345678",
		}).
		Expect().
		Status(401)

	res.JSON().Path("$.error").String().Equal(ErrInvalidCreds.Msg)
}

func TestAuthLogin_NotActivated(t *testing.T) {
	setup(t)
	defer teardown(t)

	e := heDefault(t)
	e = heHost(e, "ctfzone.test")
	e = heCSRF(e)

	res := e.POST("/api/auth/login").
		WithJSON(map[string]interface{}{
			"email":    "217@mail.com",
			"password": "217",
		}).
		Expect().
		Status(422)

	res.JSON().Path("$.error").String().Equal(ErrAccountIsNotActivated.Msg)
}

func TestAuthLogin_Success(t *testing.T) {
	setup(t)
	defer teardown(t)

	e := heDefault(t)
	e = heHost(e, "ctfzone.test")
	e = heCSRF(e)

	res := e.POST("/api/auth/login").
		WithJSON(map[string]interface{}{
			"email":    "lcbc@mail.com",
			"password": "lcbc",
		}).
		Expect().
		Status(200)

	res.Cookie("session").Value().NotEmpty()
}

func TestAuthActivate_TokenNotExists(t *testing.T) {
	setup(t)
	defer teardown(t)

	e := heDefault(t)
	e = heHost(e, "ctfzone.test")
	e = heCSRF(e)

	res := e.POST("/api/auth/activate").
		WithJSON(map[string]interface{}{
			"token": "1111111111111111111111111111111111111111111111111111111111111111",
		}).
		Expect().
		Status(404)

	res.JSON().Path("$.error").String().Equal(ErrTokenNotFound.Msg)
}

func TestAuthActivate_TokenExpired(t *testing.T) {
	setup(t)
	defer teardown(t)

	e := heDefault(t)
	e = heHost(e, "ctfzone.test")
	e = heCSRF(e)

	res := e.POST("/api/auth/activate").
		WithJSON(map[string]interface{}{
			"token": "1534e7c3fc6b3dbb411a5170b5fa94cd95f324f1ca072853ab2cb34c1378c061",
		}).
		Expect().
		Status(403)

	res.JSON().Path("$.error").String().Equal(ErrTokenIsExpired.Msg)
}

func TestAuthActivate_Success(t *testing.T) {
	setup(t)
	defer teardown(t)

	e := heDefault(t)
	e = heHost(e, "ctfzone.test")
	e = heCSRF(e)

	e.POST("/api/auth/activate").
		WithJSON(map[string]interface{}{
			"token": "5b2fae511a73355aff9a99b2164caa343889773150854d6160e4acc6c0137f17",
		}).
		Expect().
		Status(200).NoContent()
}

func TestAuthLogout_Success(t *testing.T) {
	setup(t)
	defer teardown(t)

	e := heDefault(t)
	e = heHost(e, "ctfzone.test")
	e = heCSRF(e)
	e = heAuth(e, "lcbc@mail.com", "lcbc")

	res := e.POST("/api/auth/logout").
		WithHeader("Content-Type", "application/json").
		WithText("{}").
		Expect().
		Status(200)

	res.NoContent()
}

func TestAuthSendToken_NotExistingUser(t *testing.T) {
	setup(t)
	defer teardown(t)

	e := heDefault(t)
	e = heHost(e, "ctfzone.test")
	e = heCSRF(e)

	res := e.POST("/api/auth/send-token").
		WithJSON(map[string]interface{}{
			"email": "not-existing-email@mail.com",
			"type":  "reset",
		}).
		Expect().
		Status(404)

	res.JSON().Path("$.error").String().Equal(ErrUserNotFound.Msg)
}

func TestAuthSendToken_Success(t *testing.T) {
	setup(t)
	defer teardown(t)

	e := heDefault(t)
	e = heHost(e, "ctfzone.test")
	e = heCSRF(e)

	done := make(chan struct{})

	m.On("Send",
		"reset",
		"lcbc@mail.com",
		mock.Anything).
		Return(nil).
		Run(func(mock.Arguments) { close(done) })

	res := e.POST("/api/auth/send-token").
		WithJSON(map[string]interface{}{
			"email": "lcbc@mail.com",
			"type":  "reset",
		}).
		Expect().
		Status(200)

	res.NoContent()

	select {
	case <-done:
	case <-time.After(10 * time.Millisecond):
		t.Log("Timeout")
	}

	m.AssertExpectations(t)
}

func TestAuthResetPassword_NotExistingToken(t *testing.T) {
	setup(t)
	defer teardown(t)

	e := heDefault(t)
	e = heHost(e, "ctfzone.test")
	e = heCSRF(e)

	res := e.POST("/api/auth/reset-password").
		WithJSON(map[string]interface{}{
			"token":    "1111111111111111111111111111111111111111111111111111111111111111",
			"password": "P@ssw0rd",
		}).
		Expect().
		Status(404)

	res.JSON().Path("$.error").String().Equal(ErrTokenNotFound.Msg)
}

func TestAuthResetPassword_TokenExpired(t *testing.T) {
	setup(t)
	defer teardown(t)

	e := heDefault(t)
	e = heHost(e, "ctfzone.test")
	e = heCSRF(e)

	res := e.POST("/api/auth/reset-password").
		WithJSON(map[string]interface{}{
			"token":    "150ef2dc5ab146ef23bc1a239fcb7f631006677453782dc779dc8a45a763751f",
			"password": "P@ssw0rd",
		}).
		Expect().
		Status(403)

	res.JSON().Path("$.error").String().Equal(ErrTokenIsExpired.Msg)
}

func TestAuthResetPassword_Success(t *testing.T) {
	setup(t)
	defer teardown(t)

	e := heDefault(t)
	e = heHost(e, "ctfzone.test")
	e = heCSRF(e)

	res := e.POST("/api/auth/reset-password").
		WithJSON(map[string]interface{}{
			"token":    "e2e5fbe043d0f05f37405f363a1f442124ea5de85f862e2138fcb5b9d6eff142",
			"password": "P@ssw0rd",
		}).
		Expect().
		Status(200)

	res.NoContent()
}
