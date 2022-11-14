package memory

import (
	"time"

	"github.com/company/blanksvc/pkg/models"
	"github.com/company/blanksvc/pkg/repository"
	"github.com/company/blanksvc/pkg/service/filtering"
	"github.com/company/blanksvc/pkg/service/sorting"
)

type repo struct {
	tasks []*models.Task
}

func NewRepo() repository.Repository {
	return &repo{}
}

func (r *repo) ListTasks(filter *filtering.Filter, sort *sorting.Sort, itemsOnPage int, page int) ([]*models.Task, int64, error) {
	return r.tasks, int64(len(r.tasks)), nil
}

func (r *repo) GetTask(id string) (*models.Task, error) {
	for _, task := range r.tasks {
		if task.ID == id {
			return task, nil
		}
	}
	return nil, repository.ErrTaskNotFound
}

func (r *repo) CreateTask(task *models.Task) error {
	r.tasks = append(r.tasks, task)

	return nil
}

func (r *repo) UpdateTask(taskUpdates *models.Task) error {
	for _, task := range r.tasks {
		if task.ID == taskUpdates.ID {
			task.Name = taskUpdates.Name
			task.Description = taskUpdates.Description
			task.Deadline = taskUpdates.Deadline
			task.CompletedAt = taskUpdates.CompletedAt
			updatedAt := time.Now()
			task.UpdatedAt = &updatedAt
			return nil
		}
	}

	return repository.ErrTaskNotFound
}

func (r *repo) CompleteTask(id string) error {
	for _, task := range r.tasks {
		if task.ID == id {
			completedAt := time.Now()
			task.CompletedAt = &completedAt
			return nil
		}
	}

	return repository.ErrTaskNotFound
}

func (r *repo) UncompleteTask(id string) error {
	for _, task := range r.tasks {
		if task.ID == id {
			task.CompletedAt = nil
			return nil
		}
	}

	return repository.ErrTaskNotFound
}

func (r *repo) DeleteTask(id string) error {

	return nil
}
