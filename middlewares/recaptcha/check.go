package recaptcha

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"
)

// See https://developers.google.com/recaptcha/docs/verify
var (
	checkURL  = "https://www.google.com/recaptcha/api/siteverify"
	errorText = map[string]string{
		"missing-input-secret":   "The secret parameter is missing",
		"invalid-input-secret":   "The secret parameter is invalid or malformed",
		"missing-input-response": "The response parameter is missing",
		"invalid-input-response": "The response parameter is invalid or malformed",
	}
)

type recaptchaErrors []string

func (re recaptchaErrors) Error() string {
	var err string
	for _, e := range re {
		err += errorText[e] + "; "
	}
	return err
}

type recaptchaResponse struct {
	Success   bool            `json:"success"`
	Challenge time.Time       `json:"challenge_ts"`
	Hostname  string          `json:"hostname"`
	Errors    recaptchaErrors `json:"error-codes"`
}

func check(secret, clientResponse, clientIP string) error {

	r, err := http.PostForm(checkURL, url.Values{
		"secret":   {secret},
		"response": {clientResponse},
		"remoteip": {clientIP},
	})

	if err != nil {
		return err
	}

	defer r.Body.Close()

	var res recaptchaResponse
	err = json.NewDecoder(r.Body).Decode(&res)

	if err != nil {
		return err
	}

	if !res.Success {
		return res.Errors
	}

	return nil
}
