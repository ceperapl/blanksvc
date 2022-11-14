// nolint: golint, stylecheck
package postgres

import (
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	// nolint: revive
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	bindata "github.com/golang-migrate/migrate/v4/source/go_bindata"
	"github.com/pkg/errors"
)

// Migrate applies migrations files to dsn.
func Migrate(dsn string) error {
	var err error

	s := bindata.Resource(AssetNames(), Asset)
	d, err := bindata.WithInstance(s)
	if err != nil {
		return fmt.Errorf("resource: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("go-bindata", d, dsn)
	if err != nil {
		return fmt.Errorf("instance: %w", err)
	}

	err = m.Up()
	// check on os.ErrNotExist is allowing service to rollback
	if err != nil && !(errors.Is(err, migrate.ErrNoChange) || errors.Is(err, os.ErrNotExist)) {
		return fmt.Errorf("up: %w", err)
	}
	return nil
}
