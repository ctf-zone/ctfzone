package migrations

import (
	"github.com/mattes/migrate"
	_ "github.com/mattes/migrate/database/postgres" // postgres driver for migrate
	"github.com/mattes/migrate/source/go-bindata"
	"github.com/pkg/errors"
)

func Up(dsn string) error {
	assets := bindata.Resource(AssetNames(),
		func(name string) ([]byte, error) {
			return Asset(name)
		})

	src, err := bindata.WithInstance(assets)
	if err != nil {
		return errors.Wrap(err, "model: fail to create bindata source")
	}

	mgr, err := migrate.NewWithSourceInstance("go-bindata", src, dsn)
	if err != nil {
		return errors.Wrap(err, "model: fail init migrate")
	}

	if err := mgr.Up(); err != nil && err != migrate.ErrNoChange {
		return errors.Wrap(err, "model: migrations up failed")
	}

	if err := src.Close(); err != nil {
		return errors.Wrap(err, "model: fail to close source")
	}

	if se, de := mgr.Close(); se != nil || de != nil {
		return errors.Wrap(de, "model: fail to close source")
	}

	return nil
}
