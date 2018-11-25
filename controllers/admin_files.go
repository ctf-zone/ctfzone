package controllers

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/ctf-zone/ctfzone/config"
)

func AdminFilesUpload(cfg *config.Files, baseURL string) http.HandlerFunc {

	type Response struct {
		URL string `json:"url"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseMultipartForm(cfg.MaxSize)

		file, hdr, err := r.FormFile("file")
		if err != nil {
			handleError(w, r, ErrInternal.SetError(err))
			return
		}
		defer file.Close()

		h := sha256.New()
		if _, err := io.Copy(h, file); err != nil {
			handleError(w, r, ErrInternal.SetError(err))
			return
		}

		fpath := path.Join(cfg.Path, path.Base(
			fmt.Sprintf("%s.%x", hdr.Filename, h.Sum(nil)),
		))

		file.Seek(0, os.SEEK_SET)

		f, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			handleError(w, r, ErrInternal.SetError(err))
			return
		}
		defer f.Close()

		if _, err := io.Copy(f, file); err != nil {
			handleError(w, r, ErrInternal.SetError(err))
			return
		}

		res := &Response{
			URL: fmt.Sprintf("%s/%s", baseURL, fpath),
		}

		if err := responseJSON(w, res); err != nil {
			handleError(w, r, ErrInternal.SetError(err))
		}
	}
}

func AdminFilesGetAll(c *config.Files) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		files, err := getDirFiles(c.Path)
		if err != nil {
			handleError(w, r, ErrInternal.SetError(err))
			return
		}

		if err := responseJSON(w, files); err != nil {
			handleError(w, r, ErrInternal.SetError(err))
		}
	}
}
