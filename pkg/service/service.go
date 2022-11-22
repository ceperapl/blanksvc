package service

import (
	"errors"
	"fmt"

	"github.com/company/blanksvc/pkg/common"
	"github.com/company/blanksvc/pkg/filtering"
	"github.com/company/blanksvc/pkg/models"
	"github.com/company/blanksvc/pkg/repository"
	"github.com/company/blanksvc/pkg/sorting"
	"github.com/go-kit/log"
)

type Service interface {
	ListTasks(filter *filtering.Filter, sort *sorting.Sort, itemsOnPage int, page int) ([]models.Task, int64, error)
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

func (s service) ListTasks(filter *filtering.Filter, sort *sorting.Sort, itemsOnPage int, page int) ([]models.Task, int64, error) {
	tasks, count, err := s.repo.ListTasks(filter, sort, itemsOnPage, page)
	if err != nil {
		return nil, 0, fmt.Errorf("list tasks repository: %w", err)
	}
	return tasks, count, nil
}

func (s service) GetTask(id string) (*models.Task, error) {
	task, err := s.repo.GetTask(id)
	if err != nil {
		return nil, fmt.Errorf("get task repository: %w", err)
	}
	return task, nil
}

func (s service) CreateTask(task *models.Task) (*models.Task, error) {
	if err := s.repo.CreateTask(task); err != nil {
		return nil, fmt.Errorf("create task repository: %w", err)
	}
	createdTask, err := s.GetTask(task.ID)
	if err != nil {
		return nil, fmt.Errorf("get task repository: %w", err)
	}
	return createdTask, nil
}

func (s service) UpdateTask(task *models.Task) (*models.Task, error) {
	if err := s.repo.UpdateTask(task); err != nil {
		if !errors.Is(err, common.ErrTaskNotFound) {
			return nil, fmt.Errorf("update task repository: %w", err)
		}
		if err := s.repo.CreateTask(task); err != nil {
			return nil, fmt.Errorf("create task repository: %w", err)
		}
	}
	updatedTask, err := s.GetTask(task.ID)
	if err != nil {
		return nil, fmt.Errorf("get task repository: %w", err)
	}
	return updatedTask, nil
}

func (s service) CompleteTask(id string) error {
	if err := s.repo.CompleteTask(id); err != nil {
		return fmt.Errorf("complete task repository: %w", err)
	}
	return nil
}

func (s service) UncompleteTask(id string) error {
	if err := s.repo.UncompleteTask(id); err != nil {
		return fmt.Errorf("uncomplete task repository: %w", err)
	}
	return nil
}

func (s service) DeleteTask(id string) error {
	if err := s.repo.DeleteTask(id); err != nil {
		return fmt.Errorf("delete task repository: %w", err)
	}
	return nil
}
