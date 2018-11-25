package controllers

import (
	"fmt"
	"net/http"

	"github.com/alexedwards/scs"
	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/xeipuuv/gojsonschema"

	"github.com/ctf-zone/ctfzone/controllers/schemas"
	"github.com/ctf-zone/ctfzone/middlewares/schemacheck"
)

func isLoggedIn(sm *scs.Manager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			session := sm.Load(r)

			userID, err := session.GetInt("userId")
			if err != nil || userID <= 0 {
				handleError(w, r, &Error{Code: 401, Msg: "Unauthorized request"})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func isAdmin(sm *scs.Manager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			session := sm.Load(r)

			isAdmin, err := session.GetBool("isAdmin")
			if err != nil || !isAdmin {
				handleError(w, r, &Error{Code: 403, Msg: "Access denied"})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func validate(schemaName string) func(http.Handler) http.Handler {
	return schemacheck.New(
		// Load schema from go-bindata file system.
		gojsonschema.NewReferenceLoaderFileSystem(
			fmt.Sprintf("file://%s", schemaName),
			&assetfs.AssetFS{
				Asset:     schemas.Asset,
				AssetDir:  schemas.AssetDir,
				AssetInfo: schemas.AssetInfo,
			},
		),
		schemacheck.ErrorFunc(func(w http.ResponseWriter, r *http.Request, err error) {
			t := err.(*schemacheck.Error)
			e := &Error{
				Code: 400,
				Msg:  t.Msg,
				Errs: make(map[string][]string),
			}
			for _, fe := range t.Errs {
				if _, ok := e.Errs[fe.Field]; !ok {
					e.Errs[fe.Field] = make([]string, 0)
				}
				e.Errs[fe.Field] = append(e.Errs[fe.Field], fe.Msg)
			}
			handleError(w, r, e)
		}),
	)
}

func conditional(handler func(http.Handler) http.Handler, condition bool) func(http.Handler) http.Handler {

	if condition {
		return handler
	}

	return func(next http.Handler) http.Handler {
		return next
	}
}
