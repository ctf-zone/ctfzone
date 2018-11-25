package models_test

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/go-testfixtures/testfixtures"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"upper.io/db.v3/lib/sqlbuilder"
	"upper.io/db.v3/postgresql"

	. "github.com/ctf-zone/ctfzone/models"
)

var (
	db       *Repository
	upperDB  sqlbuilder.Database
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

	sqlDB, err := sql.Open("postgres", dsn)
	if err != nil {
		return errors.Wrap(err, "fail to connect database")
	}

	upperDB, err = postgresql.New(sqlDB)
	if err != nil {
		return errors.Wrap(err, "fail init sqlbuilder.Database")
	}

	db, err = NewWithDB(upperDB)
	if err != nil {
		return errors.Wrap(err, "fail to init models")
	}

	fixtures, err = testfixtures.NewFolder(sqlDB,
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
