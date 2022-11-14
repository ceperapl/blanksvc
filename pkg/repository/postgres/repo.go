package postgres

import (
	"database/sql"
	"errors"

	"github.com/company/blanksvc/pkg/common"
	"github.com/company/blanksvc/pkg/models"
	"github.com/company/blanksvc/pkg/repository"
	"github.com/company/blanksvc/pkg/repository/postgres/transactions"
	"github.com/company/blanksvc/pkg/service/filtering"
	"github.com/company/blanksvc/pkg/service/sorting"
)

type repo struct {
	txFactory transactions.TxFactory
}

func NewRepo(txFactory transactions.TxFactory) repository.Repository {
	return &repo{
		txFactory: txFactory,
	}
}

func (r repo) ListTasks(filter *filtering.Filter, sort *sorting.Sort, itemsOnPage int, page int) ([]models.Task, int64, error) {
	sqlTx, err := r.txFactory.Begin()
	if err != nil {
		return nil, 0, err
	}
	var tasks []models.Task
	var count int64
	err = sqlTx.Do(func(sqlTx transactions.TxSQL) (errDo error) {
		query := `SELECT * from tasks`

		rows, err := sqlTx.Query(query)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			task := models.Task{}
			if errScan := rows.Scan(
				&task.ID, &task.Name, &task.Description, &task.Deadline, &task.CompletedAt, &task.CreatedAt, &task.UpdatedAt,
			); errScan != nil {
				return errScan
			}
			tasks = append(tasks, task)
		}
		err = rows.Err()
		if err != nil {
			return err
		}

		query = `SELECT count(*) from tasks`
		err = sqlTx.QueryRow(query).
			Scan(&count)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, 0, err
	}

	return tasks, count, nil
}

func (r repo) GetTask(id string) (*models.Task, error) {
	sqlTx, err := r.txFactory.Begin()
	if err != nil {
		return nil, err
	}
	var task models.Task
	err = sqlTx.Do(func(sqlTx transactions.TxSQL) (errDo error) {
		query := `SELECT * from tasks 
		WHERE id = $1`

		err := sqlTx.QueryRow(query, id).
			Scan(&task.ID, &task.Name, &task.Description, &task.Deadline, &task.CompletedAt, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, common.ErrTaskNotFound
		}
		return nil, err
	}

	return &task, nil
}

func (r repo) CreateTask(task *models.Task) error {
	sqlTx, err := r.txFactory.Begin()
	if err != nil {
		return err
	}
	err = sqlTx.Do(func(sqlTx transactions.TxSQL) (errDo error) {
		query := `INSERT INTO tasks (id, name, description, deadline) 
			VALUES($1, $2, $3, $4)`

		_, err := sqlTx.Exec(query, task.ID, task.Name, task.Description, task.Deadline)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (r repo) UpdateTask(task *models.Task) error {
	sqlTx, err := r.txFactory.Begin()
	if err != nil {
		return err
	}
	if _, getErr := r.GetTask(task.ID); err != nil {
		return getErr
	}
	err = sqlTx.Do(func(sqlTx transactions.TxSQL) (errDo error) {
		query := `UPDATE tasks SET name = $1, description = $2, deadline = $3, updated_at = now()
			WHERE id = $5`

		_, err := sqlTx.Exec(query, task.Name, task.Description, task.Deadline, task.CompletedAt, task.ID)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (r repo) CompleteTask(id string) error {
	sqlTx, err := r.txFactory.Begin()
	if err != nil {
		return err
	}
	err = sqlTx.Do(func(sqlTx transactions.TxSQL) (errDo error) {
		query := `UPDATE tasks SET completed_at = now() WHERE id = $1`

		_, err := sqlTx.Exec(query, id)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (r repo) UncompleteTask(id string) error {
	sqlTx, err := r.txFactory.Begin()
	if err != nil {
		return err
	}
	err = sqlTx.Do(func(sqlTx transactions.TxSQL) (errDo error) {
		query := `UPDATE tasks SET completed_at = null WHERE id = $1`

		_, err := sqlTx.Exec(query, id)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (r repo) DeleteTask(id string) error {
	sqlTx, err := r.txFactory.Begin()
	if err != nil {
		return err
	}
	err = sqlTx.Do(func(sqlTx transactions.TxSQL) (errDo error) {
		query := `DELETE FROM tasks WHERE id = $1`

		_, err := sqlTx.Exec(query, id)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
