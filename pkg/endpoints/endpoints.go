package endpoints

import (
	"context"
	"fmt"

	"github.com/company/blanksvc/pkg/common"
	"github.com/company/blanksvc/pkg/endpoints/middleware"
	"github.com/company/blanksvc/pkg/filtering"
	"github.com/company/blanksvc/pkg/metrics"
	"github.com/company/blanksvc/pkg/models"
	"github.com/company/blanksvc/pkg/service"
	"github.com/company/blanksvc/pkg/sorting"
	kitmetrics "github.com/go-kit/kit/metrics"
	uuid "github.com/satori/go.uuid"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
)

const (
	listTasksEndpointName      = "ListTasks"
	getTaskEndpointName        = "GetTask"
	createTaskEndpointName     = "CreateTask"
	updateTaskEndpointName     = "UpdateTask"
	completeTaskEndpointName   = "CompleteTask"
	uncompleteTaskEndpointName = "UncompleteTask"
	deleteTaskEndpointName     = "DeleteTask"
)

type ListTasksRequest struct {
	Filter      string `json:"filter"`
	Sort        string `json:"sort"`
	ItemsOnPage int    `json:"itemsOnPage"`
	Page        int    `json:"page"`
}

type ListTaskResponse struct {
	Tasks []models.Task `json:"tasks"`
	Count int64         `json:"count"`
}

type GetTaskRequest struct {
	ID string `json:"id"`
}

type CompleteTaskRequest struct {
	ID string `json:"id"`
}

type UncompleteTaskRequest struct {
	ID string `json:"id"`
}

type DeleteTaskRequest struct {
	ID string `json:"id"`
}

type Endpoints struct {
	ListTasksEndpoint      endpoint.Endpoint
	GetTaskEndpoint        endpoint.Endpoint
	CreateTaskEndpoint     endpoint.Endpoint
	UpdateTaskEndpoint     endpoint.Endpoint
	CompleteTaskEndpoint   endpoint.Endpoint
	UncompleteTaskEndpoint endpoint.Endpoint
	DeleteTaskEndpoint     endpoint.Endpoint
}

func New(svc service.Service, logger log.Logger, requestLatency kitmetrics.Histogram, requestCounter kitmetrics.Counter,
	errorCounter kitmetrics.Counter) Endpoints {
	listTasksEndpoint := middleware.LoggingMiddleware(log.With(logger, "method", listTasksEndpointName))(
		middleware.MetricsMiddleware(
			requestLatency.With(metrics.MetricsEndpointNameLabel, listTasksEndpointName),
			requestCounter.With(metrics.MetricsEndpointNameLabel, listTasksEndpointName),
			errorCounter,
		)(MakeListTasksEndpoint(svc)),
	)
	getTaskEndpoint := middleware.LoggingMiddleware(log.With(logger, "method", getTaskEndpointName))(
		middleware.MetricsMiddleware(
			requestLatency.With(metrics.MetricsEndpointNameLabel, getTaskEndpointName),
			requestCounter.With(metrics.MetricsEndpointNameLabel, getTaskEndpointName),
			errorCounter,
		)(MakeGetTaskEndpoint(svc)),
	)
	createTaskEndpoint := middleware.LoggingMiddleware(log.With(logger, "method", createTaskEndpointName))(
		middleware.MetricsMiddleware(
			requestLatency.With(metrics.MetricsEndpointNameLabel, createTaskEndpointName),
			requestCounter.With(metrics.MetricsEndpointNameLabel, createTaskEndpointName),
			errorCounter,
		)(MakeCreateTaskEndpoint(svc)),
	)
	updateTaskEndpoint := middleware.LoggingMiddleware(log.With(logger, "method", updateTaskEndpointName))(
		middleware.MetricsMiddleware(
			requestLatency.With(metrics.MetricsEndpointNameLabel, updateTaskEndpointName),
			requestCounter.With(metrics.MetricsEndpointNameLabel, updateTaskEndpointName),
			errorCounter,
		)(MakeUpdateTaskEndpoint(svc)),
	)
	completeTaskEndpoint := middleware.LoggingMiddleware(log.With(logger, "method", completeTaskEndpointName))(
		middleware.MetricsMiddleware(
			requestLatency.With(metrics.MetricsEndpointNameLabel, completeTaskEndpointName),
			requestCounter.With(metrics.MetricsEndpointNameLabel, completeTaskEndpointName),
			errorCounter,
		)(MakeCompleteTaskEndpoint(svc)),
	)
	uncompleteTaskEndpoint := middleware.LoggingMiddleware(log.With(logger, "method", uncompleteTaskEndpointName))(
		middleware.MetricsMiddleware(
			requestLatency.With(metrics.MetricsEndpointNameLabel, uncompleteTaskEndpointName),
			requestCounter.With(metrics.MetricsEndpointNameLabel, uncompleteTaskEndpointName),
			errorCounter,
		)(MakeUncompleteTaskEndpoint(svc)),
	)
	deleteTaskEndpoint := middleware.LoggingMiddleware(log.With(logger, "method", deleteTaskEndpointName))(
		middleware.MetricsMiddleware(
			requestLatency.With(metrics.MetricsEndpointNameLabel, deleteTaskEndpointName),
			requestCounter.With(metrics.MetricsEndpointNameLabel, deleteTaskEndpointName),
			errorCounter,
		)(MakeDeleteTaskEndpoint(svc)),
	)
	return Endpoints{
		ListTasksEndpoint:      listTasksEndpoint,
		GetTaskEndpoint:        getTaskEndpoint,
		CreateTaskEndpoint:     createTaskEndpoint,
		UpdateTaskEndpoint:     updateTaskEndpoint,
		CompleteTaskEndpoint:   completeTaskEndpoint,
		UncompleteTaskEndpoint: uncompleteTaskEndpoint,
		DeleteTaskEndpoint:     deleteTaskEndpoint,
	}
}

// MakeListTasksEndpoint creates ListTasks endpoint
func MakeListTasksEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		var err error
		req, ok := request.(ListTasksRequest)
		if !ok {
			return nil, fmt.Errorf("%w: ListTasksRequest", common.ErrTypeAssertion)
		}
		var filter *filtering.Filter
		if req.Filter != "" {
			filter, err = filtering.New(req.Filter, models.Task{})
			if err != nil {
				return nil, fmt.Errorf("filter creation: %w", err)
			}
		}

		var sort *sorting.Sort
		if req.Sort != "" {
			sort, err = sorting.New(req.Sort, models.Task{})
			if err != nil {
				return nil, fmt.Errorf("sort creation: %w", err)
			}
		}

		tasks, count, err := svc.ListTasks(filter, sort, req.ItemsOnPage, req.Page)
		if err != nil {
			return nil, fmt.Errorf("list tasks service: %w", err)
		}
		return ListTaskResponse{Tasks: tasks, Count: count}, nil
	}
}

// MakeGetTaskEndpoint creates GetTask endpoint
func MakeGetTaskEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(GetTaskRequest)
		if !ok {
			return nil, fmt.Errorf("%w: GetTaskRequest", common.ErrTypeAssertion)
		}
		task, err := svc.GetTask(req.ID)
		if err != nil {
			return nil, fmt.Errorf("get task service: %w", err)
		}
		return *task, nil
	}
}

// MakeCreateTaskEndpoint creates CreateTask endpoint
func MakeCreateTaskEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(*models.Task)
		if !ok {
			return nil, fmt.Errorf("%w: *Task", common.ErrTypeAssertion)
		}
		if req.ID == "" {
			req.ID = uuid.NewV4().String()
		}
		task, err := svc.CreateTask(req)
		if err != nil {
			return nil, fmt.Errorf("create task service: %w", err)
		}
		return *task, nil
	}
}

// MakeUpdateTaskEndpoint creates UpdateTask endpoint
func MakeUpdateTaskEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(*models.Task)
		if !ok {
			return nil, fmt.Errorf("%w: *Task", common.ErrTypeAssertion)
		}
		task, err := svc.UpdateTask(req)
		if err != nil {
			return nil, fmt.Errorf("update task service: %w", err)
		}
		return *task, nil
	}
}

// MakeCompleteTaskEndpoint creates CompleteTask endpoint
func MakeCompleteTaskEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(CompleteTaskRequest)
		if !ok {
			return nil, fmt.Errorf("%w: CompleteTaskRequest", common.ErrTypeAssertion)
		}
		if err := svc.CompleteTask(req.ID); err != nil {
			return nil, fmt.Errorf("complete task service: %w", err)
		}
		return nil, nil
	}
}

// MakeUncompleteTaskEndpoint creates UncompleteTask endpoint
func MakeUncompleteTaskEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(UncompleteTaskRequest)
		if !ok {
			return nil, fmt.Errorf("%w: UncompleteTaskRequest", common.ErrTypeAssertion)
		}
		if err := svc.UncompleteTask(req.ID); err != nil {
			return nil, fmt.Errorf("uncomplete task service: %w", err)
		}
		return nil, nil
	}
}

// MakeDeleteTaskEndpoint creates DeleteTask endpoint
func MakeDeleteTaskEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(DeleteTaskRequest)
		if !ok {
			return nil, fmt.Errorf("%w: DeleteTaskRequest", common.ErrTypeAssertion)
		}
		if err := svc.DeleteTask(req.ID); err != nil {
			return nil, fmt.Errorf("delete task service: %w", err)
		}
		return nil, nil
	}
}
