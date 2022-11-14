package endpoints

import (
	"context"

	"github.com/company/blanksvc/pkg/endpoints/middleware"
	"github.com/company/blanksvc/pkg/metrics"
	"github.com/company/blanksvc/pkg/models"
	"github.com/company/blanksvc/pkg/service"
	kitmetrics "github.com/go-kit/kit/metrics"

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
	Tasks []*models.Task `json:"tasks"`
	Count int64          `json:"count"`
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

func New(svc service.Service, logger log.Logger, requestLatency kitmetrics.Histogram, requestCounter kitmetrics.Counter, errorCounter kitmetrics.Counter) Endpoints {
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

func MakeListTasksEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(ListTasksRequest)
		tasks, count, err := svc.ListTasks(req.Filter, req.Sort, req.ItemsOnPage, req.Page)
		if err != nil {
			return nil, err
		}
		return ListTaskResponse{Tasks: tasks, Count: count}, nil
	}
}

func MakeGetTaskEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetTaskRequest)
		task, err := svc.GetTask(req.ID)
		if err != nil {
			return nil, err
		}
		return *task, nil
	}
}

func MakeCreateTaskEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*models.Task)
		task, err := svc.CreateTask(req)
		if err != nil {
			return nil, err
		}
		return *task, nil
	}
}

func MakeUpdateTaskEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*models.Task)
		task, err := svc.UpdateTask(req)
		if err != nil {
			return nil, err
		}
		return *task, nil
	}
}

func MakeCompleteTaskEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(CompleteTaskRequest)
		if err := svc.CompleteTask(req.ID); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

func MakeUncompleteTaskEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(UncompleteTaskRequest)
		if err := svc.UncompleteTask(req.ID); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

func MakeDeleteTaskEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(DeleteTaskRequest)
		if err := svc.DeleteTask(req.ID); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
