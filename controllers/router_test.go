package controllers_test

import (
	"bytes"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/alexedwards/scs"
	"github.com/alexedwards/scs/stores/cookiestore"
	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/gavv/httpexpect"
	"github.com/go-testfixtures/testfixtures"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/xeipuuv/gojsonschema"

	"github.com/ctf-zone/ctfzone/config"
	"github.com/ctf-zone/ctfzone/controllers/schemas"
	mailer_mock "github.com/ctf-zone/ctfzone/internal/mailer/mock"
	"github.com/ctf-zone/ctfzone/models"
	"github.com/ctf-zone/ctfzone/models/migrations"

	. "github.com/ctf-zone/ctfzone/controllers"
)

// Globals
var (
	c   config.Config
	db  *sql.DB
	dbm *models.Repository
	s   scs.Store
	m   *mailer_mock.Sender
	srv *httptest.Server
	f   *testfixtures.Context
)

// Flags
var (
	logs    bool
	verbose bool
)

func init() {
	flag.BoolVar(&logs, "logs", false, "Enables logger output.")
	flag.BoolVar(&verbose, "verbose", false, "Enables verbose HTTP printing.")
	flag.Parse()
}

func TestMain(m *testing.M) {
	if err := setupGlobals(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	ret := m.Run()

	if err := teardownGlobals(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	os.Exit(ret)
}

func setupConfig() {
	c.Game.Start = time.Now().Add(-1 * time.Hour).UTC().Format(time.RFC3339)
	c.Game.End = time.Now().Add(time.Hour).UTC().Format(time.RFC3339)
	c.Session.Lifetime = time.Hour
	c.Server.Domain = "ctfzone.test"
}

func setupMailer() {
	m = &mailer_mock.Sender{}
}

func setupDB() error {
	var err error

	dsn := os.Getenv("CTF_DB_DSN")
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		return errors.Wrap(err, "fail to connect")
	}

	dbm, err = models.New(dsn)
	if err != nil {
		return errors.Wrap(err, "fail to init models")
	}

	// Apply migrations.
	err = migrations.Up(os.Getenv("CTF_DB_DSN"))
	if err != nil {
		return errors.Wrapf(err, "fail to apply migrations, CTF_DB_DSN='%s'", os.Getenv("DB_DSN"))
	}

	// Add fixtures.
	f, err = testfixtures.NewFolder(
		db,
		&testfixtures.PostgreSQL{},
		"../models/fixtures",
	)
	if err != nil {
		return errors.Wrap(err, "fail to load fixtures")
	}

	return nil
}

func setupSessions() error {
	key, err := hex.DecodeString(strings.Repeat("11", 32))
	if err != nil {
		return err
	}
	s = cookiestore.New(key)

	return nil
}

func setupGlobals() error {
	if !logs {
		// Disable logger output.
		log.SetOutput(ioutil.Discard)
	}

	setupConfig()
	setupMailer()

	if err := setupDB(); err != nil {
		return errors.Wrap(err, "db: setup fail")
	}

	if err := setupSessions(); err != nil {
		return errors.Wrap(err, "sessions: fail to init")
	}

	routes := Router(&c, dbm, m, s)

	srv = httptest.NewServer(routes)

	return nil
}

func teardownGlobals() error {
	if db != nil {
		if err := db.Close(); err != nil {
			return errors.Wrap(err, "model: fail to close")
		}
	}

	if srv != nil {
		srv.Close()
	}

	return nil
}

func setup(t *testing.T) {
	if err := db.Ping(); err != nil {
		db, err = sql.Open("postgres", os.Getenv("CTF_DB_DSN"))
		require.NoError(t, err)
	}

	err := f.Load()
	require.NoError(t, err)
}

func teardown(t *testing.T) {}

// ===========
// = Helpers =
// ===========

func checkJSONSchema(t *testing.T, name, data string) {

	loader := gojsonschema.NewReferenceLoaderFileSystem(
		fmt.Sprintf("file:///%s", name),
		&assetfs.AssetFS{
			Asset:     schemas.Asset,
			AssetDir:  schemas.AssetDir,
			AssetInfo: schemas.AssetInfo,
		},
	)

	schema, err := gojsonschema.NewSchema(loader)
	require.NoError(t, err)

	result, err := schema.Validate(
		gojsonschema.NewStringLoader(data),
	)
	require.NoError(t, err)

	if !result.Valid() {
		fmt.Printf("Schema:\n%s\n", schemas.MustAsset(name))

		var prettyJSON bytes.Buffer
		err := json.Indent(&prettyJSON, []byte(data), "", "  ")
		require.NoError(t, err)
		fmt.Printf("Value:\n%s\n", prettyJSON.String())

		fmt.Printf("Errors:\n")
		for i, e := range result.Errors() {
			fmt.Printf("\t%d: %s\n", i, e.Description())
		}
		fmt.Printf("\n")
		t.Fail()
	}
}

func heDefault(t *testing.T) *httpexpect.Expect {
	printers := make([]httpexpect.Printer, 0)

	printers = append(printers, httpexpect.NewCompactPrinter(t))

	if verbose {
		printers = append(printers, httpexpect.NewCurlPrinter(t))
		printers = append(printers, httpexpect.NewDebugPrinter(t, true))
	}

	return httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  srv.URL,
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: printers,
	})
}

func heHost(he *httpexpect.Expect, host string) *httpexpect.Expect {
	return he.Builder(func(r *httpexpect.Request) {
		r.WithHeader("Host", host)
	})
}

func heAuth(he *httpexpect.Expect, email, password string) *httpexpect.Expect {
	return he.Builder(func(r *httpexpect.Request) {

		c := he.POST("/api/auth/login").
			WithJSON(map[string]string{
				"email":    email,
				"password": password,
			}).
			Expect().
			Status(200).
			Cookie("session").
			Value()

		c.NotEmpty()

		r.WithCookie("session", c.Raw())
	})
}

func heCSRF(he *httpexpect.Expect) *httpexpect.Expect {
	return he.Builder(func(r *httpexpect.Request) {
		r.WithHeader("X-CSRF-Token", "hVAnKbD8xC2xBgt5XGGacNTuiviEvUFmaCRDlt2Gjwo=")
		r.WithCookie("csrf-token", "hVAnKbD8xC2xBgt5XGGacNTuiviEvUFmaCRDlt2Gjwo=")
	})
}
