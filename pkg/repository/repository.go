package repository

import (
	"github.com/company/blanksvc/pkg/models"
	"github.com/company/blanksvc/pkg/service/filtering"
	"github.com/company/blanksvc/pkg/service/sorting"
)

// Repository is an abstract repository layer for data access
type Repository interface {
	ListTasks(filter *filtering.Filter, sort *sorting.Sort, itemsOnPage int, page int) ([]*models.Task, int64, error)
	GetTask(id string) (*models.Task, error)
	CreateTask(task *models.Task) error
	UpdateTask(task *models.Task) error
	CompleteTask(id string) error
	UncompleteTask(id string) error
	DeleteTask(id string) error
}
