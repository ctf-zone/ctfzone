package models_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/go-testfixtures/testfixtures"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/ctf-zone/ctfzone/models"
	. "github.com/ctf-zone/ctfzone/models"
	"github.com/ctf-zone/ctfzone/models/migrations"
)

var (
	db       *Repository
	dbx      *sqlx.DB
	fixtures *testfixtures.Context
)

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

func setupGlobals() error {
	dsn, ok := os.LookupEnv("CTF_DB_DSN")
	if !ok {
		return errors.New("empty CTF_DB_DSN")
	}

	var err error

	db, err = models.New(dsn)
	if err != nil {
		return errors.Wrap(err, "fail to init models")
	}

	dbx, err = sqlx.Connect("postgres", dsn)
	if err != nil {
		return errors.Wrap(err, "fail to connect database")
	}

	if err := migrations.Up(dsn); err != nil {
		log.Fatal(err)
	}

	fixtures, err = testfixtures.NewFolder(dbx.DB,
		&testfixtures.PostgreSQL{}, "fixtures")
	if err != nil {
		return errors.Wrap(err, "fail to load fixtures")
	}

	return nil
}

func teardownGlobals() error {
	if db != nil {
		if err := db.Close(); err != nil {
			return errors.Wrap(err, "fail to close connection to database")
		}
	}
	return nil
}

func setup(t *testing.T) {
	err := fixtures.Load()
	require.NoError(t, err)
}

func teardown(t *testing.T) {}
