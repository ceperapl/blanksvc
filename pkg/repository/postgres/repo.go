package postgres

import (
	"github.com/company/blanksvc/pkg/models"
	"github.com/company/blanksvc/pkg/repository"
	"github.com/company/blanksvc/pkg/service/filtering"
	"github.com/company/blanksvc/pkg/service/sorting"
)

type repo struct {
}

func NewRepo() repository.Repository {
	return &repo{}
}

func (r repo) ListTasks(filter *filtering.Filter, sort *sorting.Sort, itemsOnPage int, page int) ([]models.Task, int64, error) {
	return []models.Task{
		{
			ID:   "asdf",
			Name: "task1",
		},
	}, 10, nil
	// return nil, 0, nil
}

func (r repo) GetTask(id string) (*models.Task, error) {
	return &models.Task{
		ID:   "asdf",
		Name: "task1",
	}, nil
}

func (r repo) CreateTask(task *models.Task) error {
	return nil
}

func (r repo) UpdateTask(task *models.Task) error {
	return nil
}

func (r repo) CompleteTask(id string) error {
	return nil
}

func (r repo) UncompleteTask(id string) error {
	return nil
}

func (r repo) DeleteTask(id string) error {
	return nil
}
