// Package controllers
package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/alexedwards/scs"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	log "github.com/sirupsen/logrus"

	"github.com/ctf-zone/ctfzone/config"
	"github.com/ctf-zone/ctfzone/internal/mailer"
	"github.com/ctf-zone/ctfzone/middlewares/cntcheck"
	"github.com/ctf-zone/ctfzone/middlewares/csrf"
	"github.com/ctf-zone/ctfzone/middlewares/timecheck"
	"github.com/ctf-zone/ctfzone/models"
)

func init() {
	middleware.DefaultLogger = middleware.RequestLogger(
		&middleware.DefaultLogFormatter{Logger: log.StandardLogger()},
	)
}

// Router returns API HTTP handler
func Router(cfg *config.Config, db *models.Repository,
	m mailer.Sender, s scs.Store) http.Handler {

	r := chi.NewRouter()

	// CSRF
	xsrf := csrf.New(
		csrf.Header("X-CSRF-Token"),
		csrf.CookieLifetime(cfg.CSRF.Lifetime),
		csrf.CookieName("csrf-token"),
		csrf.CookieSecure(cfg.Server.TLS.Enabled),
		csrf.ErrorFunc(errorFunc(ErrInvalidCSRFToken)))

	// Session manager
	sm := scs.NewManager(s)
	sm.Persist(true)
	sm.Lifetime(cfg.Session.Lifetime)
	sm.Secure(cfg.Server.TLS.Enabled)

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(xsrf)
	r.Use(sm.Multi)

	publicRouter := PublicRouter(cfg, db, m, sm)
	adminRouter := AdminRouter(cfg, db, sm)

	r.Mount("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		host := strings.Split(requestHost(r), ":")[0]

		switch host {

		case cfg.Server.Domain:
			publicRouter.ServeHTTP(w, r)
			return

		case fmt.Sprintf("admin.%s", cfg.Server.Domain):
			adminRouter.ServeHTTP(w, r)
			return
		}

		http.Error(w, http.StatusText(404), 404)
	}))

	return r
}

func PublicRouter(cfg *config.Config, db *models.Repository,
	m mailer.Sender, sm *scs.Manager) chi.Router {

	r := chi.NewRouter()

	// Content-Type check
	cnt := cntcheck.New(
		cntcheck.ErrorFunc(errorFunc(ErrInvalidContentType)))

	// ReCaptcha
	// rc := recaptcha.New(
	// 	cfg.ReCaptcha.Secret,
	// 	recaptcha.Header(cfg.ReCaptcha.Header),
	// 	recaptcha.ErrorFunc(errorFunc(ErrInvalidCaptcha)))

	// Time check not before
	notBeforeStart := timecheck.New(cfg.Game.StartTime(), timecheck.MaxTime,
		timecheck.ErrorFunc(errorFunc(ErrCtfNotStarted)))

	// Time check only during
	during := timecheck.New(cfg.Game.StartTime(), cfg.Game.EndTime(),
		timecheck.ErrorFunc(func(w http.ResponseWriter, r *http.Request, err error) {
			switch err {
			case timecheck.ErrTooEarly:
				handleError(w, r, ErrCtfNotStarted)
			case timecheck.ErrTooLate:
				handleError(w, r, ErrCtfAlreadyEnded)
			}
		}))

	r.Route("/api", func(r chi.Router) {

		r.Use(cnt)

		r.Route("/auth", func(r chi.Router) {

			// r.With(
			// 	validate("/AuthRegister.json"),
			// 	conditional(rc, cfg.ReCaptcha.Enabled)).
			// 	Post("/register", AuthRegister(db, m))

			r.With(validate("/AuthLogin.json")).
				Post("/login", AuthLogin(db, sm))

			// r.With(validate("/AuthActivate.json")).
			// 	Post("/activate", AuthActivate(db))

			// r.With(validate("/AuthSendToken.json")).
			// 	Post("/send-token", AuthSendToken(db, m))

			// r.With(validate("/AuthResetPassword.json")).
			// 	Post("/reset-password", AuthResetPassword(db))

			r.With(isLoggedIn(sm)).Post("/logout", AuthLogout(sm))

			r.Get("/check", AuthCheck(sm))
		})

		r.Route("/user", func(r chi.Router) {
			r.Use(isLoggedIn(sm))

			r.Get("/stats", UserGetStats(cfg, db, sm))

			r.Post("/likes/{challengeId}", UserLikesCreate(db, sm))
			r.Delete("/likes/{challengeId}", UserLikesDelete(db, sm))

			r.With(validate("/UserSolutionsCreate.json"), during).
				Post("/solutions/{challengeId}", UserSolutionsCreate(&cfg.Game.Scoring, db, sm))
		})

		r.Route("/scores", func(r chi.Router) {
			r.Get("/", ScoresList(&cfg.Game.Scoring, db))
			r.Get("/ctftime", ScoresCtftimeList(&cfg.Game.Scoring, db))
		})

		r.Route("/challenges", func(r chi.Router) {
			r.Use(isLoggedIn(sm))

			r.Group(func(r chi.Router) {
				r.Use(notBeforeStart)

				r.Get("/", ChallengesList(&cfg.Game.Scoring, db, sm))
				r.Get("/{challengeId}", ChallengesGet(&cfg.Game.Scoring, db, sm))
			})
		})

		r.Route("/game", func(r chi.Router) {
			r.Get("/", GameInfo(&cfg.Game))
		})

		r.Route("/announcements", func(r chi.Router) {
			r.Get("/", AnnouncementsList(db))
		})

	})

	// Files
	filesFS := http.StripPrefix(cfg.Files.Prefix,
		http.FileServer(http.Dir(cfg.Files.Dir())))

	r.With(isLoggedIn(sm)).
		With(notBeforeStart).
		Get(cfg.Files.Prefix+"/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			filesFS.ServeHTTP(w, r)
		}))

	// Static
	// TODO: refactor
	staticDir := http.Dir(cfg.Public.StaticPath)
	staticFS := http.FileServer(staticDir)

	r.Get("/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := staticDir.Open(path.Clean(r.URL.Path))
		if os.IsNotExist(err) {
			http.ServeFile(w, r, filepath.Join(cfg.Public.StaticPath, "index.html"))
			return
		}
		staticFS.ServeHTTP(w, r)
	}))

	return r
}

func AdminRouter(cfg *config.Config, db *models.Repository, sm *scs.Manager) *chi.Mux {
	r := chi.NewRouter()

	// Content-Type check
	cnt := cntcheck.New(
		cntcheck.ErrorFunc(errorFunc(ErrInvalidContentType)))

	r.Route("/api", func(r chi.Router) {
		r.Use(cnt)

		r.Route("/auth", func(r chi.Router) {

			r.With(validate("/admin/auth/login.json")).
				Post("/login", AdminAuthLogin(cfg, db, sm))

			r.With(isAdmin(sm)).
				Post("/logout", AdminAuthLogout(sm))

			r.With(isAdmin(sm)).
				Get("/check", AdminAuthCheck())
		})

		r.Group(func(r chi.Router) {
			r.Use(isAdmin(sm))

			// Game routes.
			r.Get("/game", AdminGameInfo(&cfg.Game))

			// Users routes.
			r.Route("/users", func(r chi.Router) {
				r.Get("/", AdminUsersList(cfg, db))

				r.With(validate("/admin/users/create.json")).
					Post("/", AdminUsersCreate(db))

				r.Route("/{userId}", func(r chi.Router) {
					r.Get("/", AdminUsersGet(db))

					r.Delete("/", AdminUsersDelete(db))

					r.With(validate("/admin/users/update.json")).
						Put("/", AdminUsersUpdate(db))
				})
			})

			// Challenges routes.
			r.Route("/challenges", func(r chi.Router) {

				r.Get("/", AdminChallengesList(&cfg.Game.Scoring, db))

				r.With(validate("/admin/challenges/create.json")).
					Post("/", AdminChallengesCreate(db))

				r.Route("/{challengeId}", func(r chi.Router) {
					r.Get("/", AdminChallengesGet(&cfg.Game.Scoring, db))

					r.Delete("/", AdminChallengesDelete(db))

					r.With(validate("/admin/challenges/update.json")).
						Put("/", AdminChallengesUpdate(db))
				})
			})

			// Announcements.
			r.Route("/announcements", func(r chi.Router) {
				r.Get("/", AnnouncementsList(db))
				r.With(validate("/admin/announcements/create.json")).
					Post("/", AdminAnnouncementsCreate(db))
				r.Route("/{announcementId}", func(r chi.Router) {
					r.Get("/", AnnouncementsGet(db))
					r.With(validate("/admin/announcements/update.json")).
						Put("/", AdminAnnouncementsUpdate(db))
					r.Delete("/", AdminAnnouncementsDelete(db))
				})
			})
		})
	})

	r.Route("/api/files", func(r chi.Router) {
		r.Use(isAdmin(sm))

		r.Post("/", AdminFilesUpload(&cfg.Files, cfg.Server.BaseURL()))
		r.Get("/", AdminFilesGetAll(&cfg.Files))
	})

	// Static
	// TODO: refactor
	staticDir := http.Dir(cfg.Admin.StaticPath)
	staticFS := http.FileServer(staticDir)

	r.Get("/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := staticDir.Open(path.Clean(r.URL.Path))
		if os.IsNotExist(err) {
			http.ServeFile(w, r, filepath.Join(cfg.Admin.StaticPath, "index.html"))
			return
		}
		staticFS.ServeHTTP(w, r)
	}))

	return r
}
