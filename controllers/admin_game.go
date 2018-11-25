package controllers

import (
	"net/http"

	"github.com/ctf-zone/ctfzone/config"
)

func AdminGameInfo(c *config.Game) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := responseJSON(w, c); err != nil {
			handleError(w, r, ErrInternal.SetError(err))
			return
		}
	}
}
