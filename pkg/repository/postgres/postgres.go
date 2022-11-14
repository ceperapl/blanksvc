package postgres

import (
	"database/sql"
	"time"

	"github.com/company/blanksvc/pkg/utils"
	"github.com/go-kit/log"
)

const (
	maxOpenConns = 10
	maxIdleConns = 10
)

func New(logger log.Logger, dsn string) (*sql.DB, error) {
	// Connect to Postgres with retry
	var db *sql.DB
	if err := utils.Retry("connect to postgres DB", time.Second, 30, func() (bool, error) {
		// nolint: errcheck
		logger.Log("database", dsn)
		var errOpen error
		if db, errOpen = sql.Open("postgres", dsn); errOpen != nil {
			return false, errOpen
		}
		db.SetMaxOpenConns(maxOpenConns)
		db.SetMaxIdleConns(maxIdleConns)
		if err := db.Ping(); err != nil {
			return false, err
		}
		return true, nil
	}); err != nil {
		return nil, err
	}
	return db, nil
}
