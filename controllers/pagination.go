package controllers

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/ctf-zone/ctfzone/models"
)

func addLinkHeader(w http.ResponseWriter, baseURL string, pi *models.PagesInfo, params url.Values) {
	links := make([]string, 0)

	if !pi.Prev.IsZero() {
		v := params
		v.Del("before")
		v.Del("after")
		v.Set("before", fmt.Sprintf("%v", pi.Prev.Before))
		links = append(links, fmt.Sprintf(`<%s>; rel="prev"`, baseURL+"?"+v.Encode()))
	}

	if !pi.Next.IsZero() {
		v := params
		v.Del("before")
		v.Del("after")
		v.Set("after", fmt.Sprintf("%v", pi.Next.After))
		links = append(links, fmt.Sprintf(`<%s>; rel="next"`, baseURL+"?"+v.Encode()))
	}

	if len(links) > 0 {
		w.Header().Set("Link", strings.Join(links, ", "))
	}
}
