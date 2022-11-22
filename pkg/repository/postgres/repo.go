package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/company/blanksvc/pkg/common"
	"github.com/company/blanksvc/pkg/filtering"
	"github.com/company/blanksvc/pkg/models"
	"github.com/company/blanksvc/pkg/repository"
	"github.com/company/blanksvc/pkg/repository/postgres/transactions"
	"github.com/company/blanksvc/pkg/sorting"
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
		return nil, 0, fmt.Errorf("start transaction: %w", err)
	}
	var tasks []models.Task
	var count int64
	filterSQL, args := filter.ToPostgresSQL(1)

	err = sqlTx.Do(func(sqlTx transactions.TxSQL) error {
		query := fmt.Sprintf("SELECT * from tasks %s %s", filterSQL, sort.ToSQL())

		rows, queryErr := sqlTx.Query(query, args...)
		if queryErr != nil {
			return fmt.Errorf("query execution: %w", queryErr)
		}
		defer rows.Close()

		for rows.Next() {
			task := models.Task{}
			if errScan := rows.Scan(
				&task.ID, &task.Name, &task.Description, &task.Deadline, &task.CompletedAt, &task.CreatedAt, &task.UpdatedAt,
			); errScan != nil {
				return fmt.Errorf("scan rows: %w", errScan)
			}
			tasks = append(tasks, task)
		}
		rowsErr := rows.Err()
		if rowsErr != nil {
			// nolint: wrapcheck
			return rowsErr
		}

		query = fmt.Sprintf("SELECT count(*) from tasks %s %s", filterSQL, sort.ToSQL())
		queryErr = sqlTx.QueryRow(query, args...).
			Scan(&count)
		if queryErr != nil {
			return fmt.Errorf("query execution: %w", queryErr)
		}
		return nil
	})

	if err != nil {
		return nil, 0, fmt.Errorf("doing transaction: %w", err)
	}

	return tasks, count, nil
}

func (r repo) GetTask(id string) (*models.Task, error) {
	sqlTx, err := r.txFactory.Begin()
	if err != nil {
		return nil, fmt.Errorf("start transaction: %w", err)
	}
	var task models.Task
	err = sqlTx.Do(func(sqlTx transactions.TxSQL) error {
		query := `SELECT * from tasks 
		WHERE id = $1`

		queryErr := sqlTx.QueryRow(query, id).
			Scan(&task.ID, &task.Name, &task.Description, &task.Deadline, &task.CompletedAt, &task.CreatedAt, &task.UpdatedAt)
		if queryErr != nil {
			return fmt.Errorf("query execution: %w", queryErr)
		}
		return nil
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, common.ErrTaskNotFound
		}
		return nil, fmt.Errorf("doing transaction: %w", err)
	}

	return &task, nil
}

func (r repo) CreateTask(task *models.Task) error {
	sqlTx, err := r.txFactory.Begin()
	if err != nil {
		return fmt.Errorf("start transaction: %w", err)
	}
	err = sqlTx.Do(func(sqlTx transactions.TxSQL) error {
		query := `INSERT INTO tasks (id, name, description, deadline) 
			VALUES($1, $2, $3, $4)`

		_, execErr := sqlTx.Exec(query, task.ID, task.Name, task.Description, task.Deadline)
		if execErr != nil {
			return fmt.Errorf("exec run: %w", execErr)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("doing transaction: %w", err)
	}

	return nil
}

func (r repo) UpdateTask(task *models.Task) error {
	sqlTx, err := r.txFactory.Begin()
	if err != nil {
		return fmt.Errorf("start transaction: %w", err)
	}
	if _, getErr := r.GetTask(task.ID); err != nil {
		return getErr
	}
	err = sqlTx.Do(func(sqlTx transactions.TxSQL) error {
		query := `UPDATE tasks SET name = $1, description = $2, deadline = $3, updated_at = now()
			WHERE id = $5`

		_, execErr := sqlTx.Exec(query, task.Name, task.Description, task.Deadline, task.CompletedAt, task.ID)
		if execErr != nil {
			return fmt.Errorf("exec run: %w", execErr)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("doing transaction: %w", err)
	}

	return nil
}

func (r repo) CompleteTask(id string) error {
	sqlTx, err := r.txFactory.Begin()
	if err != nil {
		return fmt.Errorf("start transaction: %w", err)
	}
	err = sqlTx.Do(func(sqlTx transactions.TxSQL) error {
		query := `UPDATE tasks SET completed_at = now() WHERE id = $1`

		_, execErr := sqlTx.Exec(query, id)
		if execErr != nil {
			return fmt.Errorf("exec run: %w", execErr)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("doing transaction: %w", err)
	}

	return nil
}

func (r repo) UncompleteTask(id string) error {
	sqlTx, err := r.txFactory.Begin()
	if err != nil {
		return fmt.Errorf("start transaction: %w", err)
	}
	err = sqlTx.Do(func(sqlTx transactions.TxSQL) error {
		query := `UPDATE tasks SET completed_at = null WHERE id = $1`

		_, execErr := sqlTx.Exec(query, id)
		if execErr != nil {
			return fmt.Errorf("exec run: %w", execErr)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("doing transaction: %w", err)
	}

	return nil
}

func (r repo) DeleteTask(id string) error {
	sqlTx, err := r.txFactory.Begin()
	if err != nil {
		return fmt.Errorf("start transaction: %w", err)
	}
	err = sqlTx.Do(func(sqlTx transactions.TxSQL) error {
		query := `DELETE FROM tasks WHERE id = $1`

		_, execErr := sqlTx.Exec(query, id)
		if execErr != nil {
			return fmt.Errorf("exec run: %w", execErr)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("doing transaction: %w", err)
	}

	return nil
}
