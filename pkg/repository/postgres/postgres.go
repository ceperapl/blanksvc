package postgres

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/company/blanksvc/pkg/utils"
	"github.com/go-kit/log"
)

const (
	maxOpenConns = 10
	maxIdleConns = 10

	retryInterval    = time.Second
	retryMaxAttempts = 30
)

func New(logger log.Logger, dsn string) (*sql.DB, error) {
	// Connect to Postgres with retry
	var db *sql.DB
	if err := utils.Retry("connect to the postgres DB", retryInterval, retryMaxAttempts, func() (bool, error) {
		// nolint: errcheck
		logger.Log("database", dsn)
		var errOpen error
		if db, errOpen = sql.Open("postgres", dsn); errOpen != nil {
			return false, fmt.Errorf("open database: %w", errOpen)
		}
		db.SetMaxOpenConns(maxOpenConns)
		db.SetMaxIdleConns(maxIdleConns)
		if pingErr := db.Ping(); pingErr != nil {
			return false, fmt.Errorf("ping db: %w", pingErr)
		}
		return true, nil
	}); err != nil {
		return nil, fmt.Errorf("retry connecting to the postgres db: %w", err)
	}
	return db, nil
}
