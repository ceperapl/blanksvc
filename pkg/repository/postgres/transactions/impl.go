package transactions

import "log"

type txImpl struct {
	TxSQL
}

// NewTx makes transaction implementation that handles transactions steps
func NewTx(sqlTx TxSQL) Tx {
	return &txImpl{TxSQL: sqlTx}
}

func (tx *txImpl) Do(fn TxHandler) (err error) {
	err = fn(tx)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			log.Printf("rollback failed: %v", rollbackErr)
		}
	} else {
		err = tx.Commit()
	}
	return err
}
