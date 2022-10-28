package service

import (
	"fmt"

	"github.com/company/blanksvc/pkg/models"
	"github.com/company/blanksvc/pkg/repository"
	"github.com/company/blanksvc/pkg/service/filtering"
	"github.com/company/blanksvc/pkg/service/sorting"
	"github.com/go-kit/log"
	uuid "github.com/satori/go.uuid"
)

type Service interface {
	Hello(name string) (string, error)
	ListTasks(filter string, sort string, itemsOnPage int, page int) ([]models.Task, int64, error)
	GetTask(id string) (*models.Task, error)
	CreateTask(task *models.Task) (*models.Task, error)
	UpdateTask(task *models.Task) (*models.Task, error)
	CompleteTask(id string) error
	UncompleteTask(id string) error
	DeleteTask(id string) error
}

type service struct {
	logger log.Logger
	repo   repository.Repository
}

func New(logger log.Logger, repo repository.Repository) Service {
	return service{
		logger: logger,
		repo:   repo,
	}
}

func (s service) Hello(name string) (string, error) {
	if name == "" {
		return "", ErrEmptyString
	}

	return fmt.Sprintf("Hello %s", name), nil
}

func (s service) ListTasks(filter string, sort string, itemsOnPage int, page int) ([]models.Task, int64, error) {
	taskFilter, err := filtering.NewFilter(filter, models.Task{})
	if err != nil {
		return nil, 0, err
	}
	taskSort, err := sorting.NewSort(sort, models.Task{})
	if err != nil {
		return nil, 0, err
	}
	return s.repo.ListTasks(taskFilter, taskSort, itemsOnPage, page)
}

func (s service) GetTask(id string) (*models.Task, error) {
	return s.repo.GetTask(id)
}

func (s service) CreateTask(task *models.Task) (*models.Task, error) {
	if task.ID == "" {
		task.ID = uuid.NewV4().String()
	}
	if err := s.repo.CreateTask(task); err != nil {
		return nil, err
	}
	return s.repo.GetTask(task.ID)
}

func (s service) UpdateTask(task *models.Task) (*models.Task, error) {
	if err := s.repo.UpdateTask(task); err != nil {
		return nil, err
	}
	return s.repo.GetTask(task.ID)
}

func (s service) CompleteTask(id string) error {
	return s.repo.CompleteTask(id)
}

func (s service) UncompleteTask(id string) error {
	return s.repo.UncompleteTask(id)
}

func (s service) DeleteTask(id string) error {
	return s.repo.DeleteTask(id)
}
