package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// Error represents an error that can be marshaled to json and sended to user.
type Error struct {
	Err  error               `json:"-"`
	Code int                 `json:"-"`
	Msg  string              `json:"error"`
	Errs map[string][]string `json:"errors,omitempty"`
}

// Allow Error to satisfy error interface.
func (e *Error) Error() string {
	var s string
	for field, err := range e.Errs {
		s += fmt.Sprintf("%s: %s;\n", field, err)
	}
	return s
}

func (e *Error) SetFieldsErrors(errs map[string][]string) *Error {
	e.Errs = errs
	return e
}

func (e *Error) SetError(err error) *Error {
	e.Err = err
	return e
}

func (e *Error) SetMessage(msg string) *Error {
	e.Msg = msg
	return e
}

func handleError(w http.ResponseWriter, r *http.Request, e *Error) {

	if e.Code == 500 {
		log.Error(e.Err)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(e.Code)
	if err := json.NewEncoder(w).Encode(e); err != nil {
		log.Error(err)
	}
}

func errorFunc(e *Error) func(http.ResponseWriter, *http.Request, error) {
	return func(w http.ResponseWriter, r *http.Request, err error) {
		handleError(w, r, e.SetError(err))
	}
}

var (
	ErrInvalidQueryParams    = &Error{Code: 400, Msg: "Invalid query params"}
	ErrInvalidJSON           = &Error{Code: 400, Msg: "Invalid JSON"}
	ErrInvalidID             = &Error{Code: 400, Msg: "Invalid ID supplied"}
	ErrInvalidContentType    = &Error{Code: 400, Msg: "Invalid Content-Type"}
	ErrUnauthorizedRequest   = &Error{Code: 401, Msg: "Unauthorized request"}
	ErrInvalidCreds          = &Error{Code: 401, Msg: "Invalid credentials"}
	ErrTokenIsExpired        = &Error{Code: 403, Msg: "Token is expired"}
	ErrAccessDenied          = &Error{Code: 403, Msg: "Access denied"}
	ErrInvalidCSRFToken      = &Error{Code: 403, Msg: "Invalid CSRF token"}
	ErrInvalidCaptcha        = &Error{Code: 403, Msg: "Invalid captcha"}
	ErrCtfNotStarted         = &Error{Code: 403, Msg: "Соревнование еще не началось"}
	ErrCtfAlreadyEnded       = &Error{Code: 403, Msg: "Соревнование закончилось"}
	ErrPathNotFound          = &Error{Code: 404, Msg: "Path not found"}
	ErrTokenNotFound         = &Error{Code: 404, Msg: "Token not found"}
	ErrChallengeNotFound     = &Error{Code: 404, Msg: "Challenge not found"}
	ErrAnnouncementNotFound  = &Error{Code: 404, Msg: "Announcement not found"}
	ErrUserNotFound          = &Error{Code: 404, Msg: "User not found"}
	ErrDuplicate             = &Error{Code: 409, Msg: "Duplicate entry"}
	ErrConditionFailed       = &Error{Code: 412, Msg: "Необходимо сначала решить зависимое задание"}
	ErrInvalidFlag           = &Error{Code: 418, Msg: "Invalid flag"}
	ErrAccountIsNotActivated = &Error{Code: 422, Msg: "Account is not activated"}
	ErrInternal              = &Error{Code: 500, Msg: "Internal error"}
)
