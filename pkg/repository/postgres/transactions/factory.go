package transactions

import (
	"database/sql"
)

// TxFactory - transactions factory
type TxFactory interface {
	Begin() (Tx, error)
}

// TxSQL
type TxSQL interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Rollback() error
	Commit() error
}

// TxHandler
type TxHandler func(tx TxSQL) error

// Tx
type Tx interface {
	TxSQL
	Do(fn TxHandler) error
}

type txFactory struct {
	db *sql.DB
}

// New constructs new TxFactory
func New(db *sql.DB) TxFactory {
	return &txFactory{db}
}

func (f *txFactory) Begin() (Tx, error) {
	tx, err := f.db.Begin()
	if err != nil {
		return nil, err
	}

	return NewTx(tx), err
}
