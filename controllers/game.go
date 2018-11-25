package controllers

import (
	"net/http"
	"time"

	"github.com/ctf-zone/ctfzone/config"
)

// GameInfo return information about CTF
// summary: Returns information about game
// tags: [game]
func GameInfo(c *config.Game) http.HandlerFunc {

	type Response struct {
		Status string    `json:"status"`
		Start  time.Time `json:"start"`
		End    time.Time `json:"end"`
	}

	start, end := c.StartTime(), c.EndTime()

	return func(w http.ResponseWriter, r *http.Request) {

		info := &Response{
			Start: start.UTC(),
			End:   end.UTC(),
		}

		now := time.Now()

		if now.Before(start) {
			info.Status = "countdown"
		} else if now.After(start) && now.Before(end) {
			info.Status = "started"
		} else if now.After(end) {
			info.Status = "ended"
		}

		// schema: GameInfo
		if err := responseJSON(w, info); err != nil {
			handleError(w, r, ErrInternal.SetError(err))
			return
		}
	}
}
