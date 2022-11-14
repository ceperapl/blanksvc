package service

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/company/blanksvc/pkg/common"
	"github.com/company/blanksvc/pkg/models"
	"github.com/company/blanksvc/pkg/repository/mocks"
	"github.com/company/blanksvc/pkg/service/filtering"
	"github.com/company/blanksvc/pkg/service/sorting"
	"github.com/go-kit/log"
	"github.com/stretchr/testify/assert"
)

func TestListTasks(t *testing.T) {
	t.Parallel()

	t.Run("successful case", func(t *testing.T) {
		t.Parallel()
		repoMock := mocks.Repository{}
		now := time.Now()
		description := "test-description"
		tasks := []models.Task{
			{
				ID:          "4ed34c91-5f6c-4316-82bf-72a211003e21",
				Name:        "test-task3",
				Description: &description,
				Deadline:    "2022-01-01",
				CompletedAt: nil,
				CreatedAt:   now,
				UpdatedAt:   nil,
			},
			{
				ID:          "4ed34c91-5f6c-4316-82bf-72a211003e22",
				Name:        "test-task2",
				Description: &description,
				Deadline:    "2022-01-01",
				CompletedAt: &now,
				CreatedAt:   now,
				UpdatedAt:   nil,
			},
			{
				ID:          "4ed34c91-5f6c-4316-82bf-72a211003e23",
				Name:        "test-task1",
				Description: &description,
				Deadline:    "2022-01-01",
				CompletedAt: nil,
				CreatedAt:   now,
				UpdatedAt:   &now,
			},
		}
		var filter *filtering.Filter
		var sort *sorting.Sort
		itemsOnPage := 0
		page := 1
		repoMock.On("ListTasks", filter, sort, itemsOnPage, page).Return(
			func(filter *filtering.Filter, sort *sorting.Sort, itemsOnPage int, page int) []models.Task {
				return tasks
			},
			func(filter *filtering.Filter, sort *sorting.Sort, itemsOnPage int, page int) int64 {
				return int64(len(tasks))
			},
			func(filter *filtering.Filter, sort *sorting.Sort, itemsOnPage int, page int) error { return nil },
		)
		logger := log.NewLogfmtLogger(os.Stderr)
		service := New(logger, &repoMock)
		listedTasks, count, err := service.ListTasks(nil, nil, itemsOnPage, page)
		assert.NoError(t, err)
		assert.Equal(t, int64(len(tasks)), count)
		assert.Equal(t, tasks, listedTasks)
	})

	t.Run("success case - empty list", func(t *testing.T) {
		t.Parallel()
		repoMock := mocks.Repository{}
		tasks := []models.Task{}
		var filter *filtering.Filter
		var sort *sorting.Sort
		itemsOnPage := 0
		page := 1
		repoMock.On("ListTasks", filter, sort, itemsOnPage, page).Return(
			func(filter *filtering.Filter, sort *sorting.Sort, itemsOnPage int, page int) []models.Task {
				return tasks
			},
			func(filter *filtering.Filter, sort *sorting.Sort, itemsOnPage int, page int) int64 {
				return int64(len(tasks))
			},
			func(filter *filtering.Filter, sort *sorting.Sort, itemsOnPage int, page int) error { return nil },
		)
		logger := log.NewLogfmtLogger(os.Stderr)
		service := New(logger, &repoMock)
		listedTasks, count, err := service.ListTasks(nil, nil, itemsOnPage, page)
		assert.NoError(t, err)
		assert.Equal(t, int64(len(tasks)), count)
		assert.Equal(t, tasks, listedTasks)
	})

	t.Run("failed to list tasks", func(t *testing.T) {
		t.Parallel()
		repoMock := mocks.Repository{}
		tasks := []models.Task{}
		var filter *filtering.Filter
		var sort *sorting.Sort
		itemsOnPage := 0
		page := 1
		expectedError := errors.New("failed to list tasks")
		repoMock.On("ListTasks", filter, sort, itemsOnPage, page).Return(
			func(filter *filtering.Filter, sort *sorting.Sort, itemsOnPage int, page int) []models.Task {
				return tasks
			},
			func(filter *filtering.Filter, sort *sorting.Sort, itemsOnPage int, page int) int64 {
				return int64(len(tasks))
			},
			func(filter *filtering.Filter, sort *sorting.Sort, itemsOnPage int, page int) error {
				return expectedError
			},
		)
		logger := log.NewLogfmtLogger(os.Stderr)
		service := New(logger, &repoMock)
		_, _, err := service.ListTasks(nil, nil, itemsOnPage, page)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, expectedError))
	})
}

func TestGetTask(t *testing.T) {
	t.Parallel()

	t.Run("successful case", func(t *testing.T) {
		t.Parallel()
		repoMock := mocks.Repository{}
		now := time.Now()
		description := "test-description"
		task := models.Task{
			ID:          "4ed34c91-5f6c-4316-82bf-72a211003e21",
			Name:        "test-task3",
			Description: &description,
			Deadline:    "2022-01-01",
			CompletedAt: nil,
			CreatedAt:   now,
			UpdatedAt:   nil,
		}
		repoMock.On("GetTask", task.ID).Return(
			func(id string) *models.Task {
				return &task
			},
			func(id string) error { return nil },
		)
		logger := log.NewLogfmtLogger(os.Stderr)
		service := New(logger, &repoMock)
		receivedTask, err := service.GetTask(task.ID)
		assert.NoError(t, err)
		assert.Equal(t, &task, receivedTask)
	})

	t.Run("failed to get task", func(t *testing.T) {
		t.Parallel()
		repoMock := mocks.Repository{}
		now := time.Now()
		description := "test-description"
		task := models.Task{
			ID:          "4ed34c91-5f6c-4316-82bf-72a211003e21",
			Name:        "test-task3",
			Description: &description,
			Deadline:    "2022-01-01",
			CompletedAt: nil,
			CreatedAt:   now,
			UpdatedAt:   nil,
		}
		expectedError := errors.New("failed to get task")
		repoMock.On("GetTask", task.ID).Return(
			func(id string) *models.Task {
				return &task
			},
			func(id string) error {
				return expectedError
			},
		)
		logger := log.NewLogfmtLogger(os.Stderr)
		service := New(logger, &repoMock)
		_, err := service.GetTask(task.ID)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, expectedError))
	})
}

func TestCreateTask(t *testing.T) {
	t.Parallel()

	t.Run("successful case", func(t *testing.T) {
		t.Parallel()
		repoMock := mocks.Repository{}
		now := time.Now()
		description := "test-description"
		task := models.Task{
			ID:          "4ed34c91-5f6c-4316-82bf-72a211003e21",
			Name:        "test-task3",
			Description: &description,
			Deadline:    "2022-01-01",
			CompletedAt: nil,
			CreatedAt:   now,
			UpdatedAt:   nil,
		}
		repoMock.On("CreateTask", &task).Return(
			func(task *models.Task) error { return nil },
		)
		repoMock.On("GetTask", task.ID).Return(
			func(id string) *models.Task {
				return &task
			},
			func(id string) error { return nil },
		)
		logger := log.NewLogfmtLogger(os.Stderr)
		service := New(logger, &repoMock)
		createdTask, err := service.CreateTask(&task)
		assert.NoError(t, err)
		assert.Equal(t, &task, createdTask)
	})

	t.Run("failed to create task", func(t *testing.T) {
		t.Parallel()
		repoMock := mocks.Repository{}
		now := time.Now()
		description := "test-description"
		task := models.Task{
			ID:          "4ed34c91-5f6c-4316-82bf-72a211003e21",
			Name:        "test-task3",
			Description: &description,
			Deadline:    "2022-01-01",
			CompletedAt: nil,
			CreatedAt:   now,
			UpdatedAt:   nil,
		}
		expectedError := errors.New("failed to create task")
		repoMock.On("CreateTask", &task).Return(
			func(task *models.Task) error {
				return expectedError
			},
		)
		logger := log.NewLogfmtLogger(os.Stderr)
		service := New(logger, &repoMock)
		_, err := service.CreateTask(&task)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, expectedError))
	})

	t.Run("failed to get task", func(t *testing.T) {
		t.Parallel()
		repoMock := mocks.Repository{}
		now := time.Now()
		description := "test-description"
		task := models.Task{
			ID:          "4ed34c91-5f6c-4316-82bf-72a211003e21",
			Name:        "test-task3",
			Description: &description,
			Deadline:    "2022-01-01",
			CompletedAt: nil,
			CreatedAt:   now,
			UpdatedAt:   nil,
		}
		repoMock.On("CreateTask", &task).Return(
			func(task *models.Task) error { return nil },
		)
		expectedError := errors.New("failed to get task")
		repoMock.On("GetTask", task.ID).Return(
			func(id string) *models.Task {
				return &task
			},
			func(id string) error {
				return expectedError
			},
		)
		logger := log.NewLogfmtLogger(os.Stderr)
		service := New(logger, &repoMock)
		_, err := service.CreateTask(&task)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, expectedError))
	})
}

func TestUpdateTask(t *testing.T) {
	t.Parallel()

	t.Run("successful case", func(t *testing.T) {
		t.Parallel()
		repoMock := mocks.Repository{}
		now := time.Now()
		description := "test-description"
		task := models.Task{
			ID:          "4ed34c91-5f6c-4316-82bf-72a211003e21",
			Name:        "test-task3",
			Description: &description,
			Deadline:    "2022-01-01",
			CompletedAt: nil,
			CreatedAt:   now,
			UpdatedAt:   nil,
		}
		repoMock.On("UpdateTask", &task).Return(
			func(task *models.Task) error { return nil },
		)
		repoMock.On("GetTask", task.ID).Return(
			func(id string) *models.Task {
				return &task
			},
			func(id string) error { return nil },
		)
		logger := log.NewLogfmtLogger(os.Stderr)
		service := New(logger, &repoMock)
		updatedTask, err := service.UpdateTask(&task)
		assert.NoError(t, err)
		assert.Equal(t, &task, updatedTask)
	})

	t.Run("successful case - task is not found", func(t *testing.T) {
		t.Parallel()
		repoMock := mocks.Repository{}
		now := time.Now()
		description := "test-description"
		task := models.Task{
			ID:          "4ed34c91-5f6c-4316-82bf-72a211003e21",
			Name:        "test-task3",
			Description: &description,
			Deadline:    "2022-01-01",
			CompletedAt: nil,
			CreatedAt:   now,
			UpdatedAt:   nil,
		}
		repoMock.On("UpdateTask", &task).Return(
			func(task *models.Task) error {
				return common.ErrTaskNotFound
			},
		)
		repoMock.On("CreateTask", &task).Return(
			func(task *models.Task) error { return nil },
		)
		repoMock.On("GetTask", task.ID).Return(
			func(id string) *models.Task {
				return &task
			},
			func(id string) error { return nil },
		)
		logger := log.NewLogfmtLogger(os.Stderr)
		service := New(logger, &repoMock)
		updatedTask, err := service.UpdateTask(&task)
		assert.NoError(t, err)
		assert.Equal(t, &task, updatedTask)
	})

	t.Run("failed to create task", func(t *testing.T) {
		t.Parallel()
		repoMock := mocks.Repository{}
		now := time.Now()
		description := "test-description"
		task := models.Task{
			ID:          "4ed34c91-5f6c-4316-82bf-72a211003e21",
			Name:        "test-task3",
			Description: &description,
			Deadline:    "2022-01-01",
			CompletedAt: nil,
			CreatedAt:   now,
			UpdatedAt:   nil,
		}
		repoMock.On("UpdateTask", &task).Return(
			func(task *models.Task) error {
				return common.ErrTaskNotFound
			},
		)
		expectedError := errors.New("failed to create task")
		repoMock.On("CreateTask", &task).Return(
			func(task *models.Task) error {
				return expectedError
			},
		)
		logger := log.NewLogfmtLogger(os.Stderr)
		service := New(logger, &repoMock)
		_, err := service.UpdateTask(&task)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, expectedError))
	})

	t.Run("failed to update task", func(t *testing.T) {
		t.Parallel()
		repoMock := mocks.Repository{}
		now := time.Now()
		description := "test-description"
		task := models.Task{
			ID:          "4ed34c91-5f6c-4316-82bf-72a211003e21",
			Name:        "test-task3",
			Description: &description,
			Deadline:    "2022-01-01",
			CompletedAt: nil,
			CreatedAt:   now,
			UpdatedAt:   nil,
		}
		expectedError := errors.New("failed to update task")
		repoMock.On("UpdateTask", &task).Return(
			func(task *models.Task) error {
				return expectedError
			},
		)
		logger := log.NewLogfmtLogger(os.Stderr)
		service := New(logger, &repoMock)
		_, err := service.UpdateTask(&task)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, expectedError))
	})

	t.Run("failed to get task", func(t *testing.T) {
		t.Parallel()
		repoMock := mocks.Repository{}
		now := time.Now()
		description := "test-description"
		task := models.Task{
			ID:          "4ed34c91-5f6c-4316-82bf-72a211003e21",
			Name:        "test-task3",
			Description: &description,
			Deadline:    "2022-01-01",
			CompletedAt: nil,
			CreatedAt:   now,
			UpdatedAt:   nil,
		}
		repoMock.On("UpdateTask", &task).Return(
			func(task *models.Task) error { return nil },
		)
		expectedError := errors.New("failed to get task")
		repoMock.On("GetTask", task.ID).Return(
			func(id string) *models.Task {
				return &task
			},
			func(id string) error {
				return expectedError
			},
		)
		logger := log.NewLogfmtLogger(os.Stderr)
		service := New(logger, &repoMock)
		_, err := service.UpdateTask(&task)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, expectedError))
	})
}

func TestCompleteTask(t *testing.T) {
	t.Parallel()

	t.Run("successful case", func(t *testing.T) {
		t.Parallel()
		repoMock := mocks.Repository{}
		taskID := "4ed34c91-5f6c-4316-82bf-72a211003e21"
		repoMock.On("CompleteTask", taskID).Return(
			func(id string) error { return nil },
		)
		logger := log.NewLogfmtLogger(os.Stderr)
		service := New(logger, &repoMock)
		err := service.CompleteTask(taskID)
		assert.NoError(t, err)
	})

	t.Run("failed to complete task", func(t *testing.T) {
		t.Parallel()
		repoMock := mocks.Repository{}
		taskID := "4ed34c91-5f6c-4316-82bf-72a211003e21"
		expectedError := errors.New("failed to complete task")
		repoMock.On("CompleteTask", taskID).Return(
			func(id string) error {
				return expectedError
			},
		)
		logger := log.NewLogfmtLogger(os.Stderr)
		service := New(logger, &repoMock)
		err := service.CompleteTask(taskID)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, expectedError))
	})
}

func TestUncompleteTask(t *testing.T) {
	t.Parallel()

	t.Run("successful case", func(t *testing.T) {
		t.Parallel()
		repoMock := mocks.Repository{}
		taskID := "4ed34c91-5f6c-4316-82bf-72a211003e21"
		repoMock.On("UncompleteTask", taskID).Return(
			func(id string) error { return nil },
		)
		logger := log.NewLogfmtLogger(os.Stderr)
		service := New(logger, &repoMock)
		err := service.UncompleteTask(taskID)
		assert.NoError(t, err)
	})

	t.Run("failed to complete task", func(t *testing.T) {
		t.Parallel()
		repoMock := mocks.Repository{}
		taskID := "4ed34c91-5f6c-4316-82bf-72a211003e21"
		expectedError := errors.New("failed to uncomplete task")
		repoMock.On("UncompleteTask", taskID).Return(
			func(id string) error {
				return expectedError
			},
		)
		logger := log.NewLogfmtLogger(os.Stderr)
		service := New(logger, &repoMock)
		err := service.UncompleteTask(taskID)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, expectedError))
	})
}

func TestDeleteTask(t *testing.T) {
	t.Parallel()

	t.Run("successful case", func(t *testing.T) {
		t.Parallel()
		repoMock := mocks.Repository{}
		taskID := "4ed34c91-5f6c-4316-82bf-72a211003e21"
		repoMock.On("DeleteTask", taskID).Return(
			func(id string) error { return nil },
		)
		logger := log.NewLogfmtLogger(os.Stderr)
		service := New(logger, &repoMock)
		err := service.DeleteTask(taskID)
		assert.NoError(t, err)
	})

	t.Run("failed to complete task", func(t *testing.T) {
		t.Parallel()
		repoMock := mocks.Repository{}
		taskID := "4ed34c91-5f6c-4316-82bf-72a211003e21"
		expectedError := errors.New("failed to delete task")
		repoMock.On("DeleteTask", taskID).Return(
			func(id string) error {
				return expectedError
			},
		)
		logger := log.NewLogfmtLogger(os.Stderr)
		service := New(logger, &repoMock)
		err := service.DeleteTask(taskID)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, expectedError))
	})
}
