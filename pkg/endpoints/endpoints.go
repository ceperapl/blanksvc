package endpoints

import (
	"context"

	"github.com/company/blanksvc/pkg/endpoints/middleware"
	"github.com/company/blanksvc/pkg/models"
	"github.com/company/blanksvc/pkg/service"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
)

type HelloRequest struct {
	Name string `json:"name"`
}

type HelloResponse struct {
	Greeting string `json:"greeting"`
}

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
	HelloEndpoint          endpoint.Endpoint
	ListTasksEndpoint      endpoint.Endpoint
	GetTaskEndpoint        endpoint.Endpoint
	CreateTaskEndpoint     endpoint.Endpoint
	UpdateTaskEndpoint     endpoint.Endpoint
	CompleteTaskEndpoint   endpoint.Endpoint
	UncompleteTaskEndpoint endpoint.Endpoint
	DeleteTaskEndpoint     endpoint.Endpoint
}

func New(svc service.Service, logger log.Logger, requestLatencyMetric *prometheus.HistogramVec) Endpoints {
	helloEndpoint := middleware.LoggingMiddleware(log.With(logger, "method", "Hello"))(
		middleware.MetricsMiddleware(requestLatencyMetric)(
			MakeHelloEndpoint(svc),
		),
	)
	listTasksEndpoint := middleware.LoggingMiddleware(log.With(logger, "method", "ListTasks"))(
		middleware.MetricsMiddleware(requestLatencyMetric)(
			MakeListTasksEndpoint(svc),
		),
	)
	getTaskEndpoint := middleware.LoggingMiddleware(log.With(logger, "method", "GetTask"))(
		middleware.MetricsMiddleware(requestLatencyMetric)(
			MakeGetTaskEndpoint(svc),
		),
	)
	createTaskEndpoint := middleware.LoggingMiddleware(log.With(logger, "method", "CreateTask"))(
		middleware.MetricsMiddleware(requestLatencyMetric)(
			MakeCreateTaskEndpoint(svc),
		),
	)
	updateTaskEndpoint := middleware.LoggingMiddleware(log.With(logger, "method", "UpdateTask"))(
		middleware.MetricsMiddleware(requestLatencyMetric)(
			MakeUpdateTaskEndpoint(svc),
		),
	)
	completeTaskEndpoint := middleware.LoggingMiddleware(log.With(logger, "method", "CompleteTask"))(
		middleware.MetricsMiddleware(requestLatencyMetric)(
			MakeCompleteTaskEndpoint(svc),
		),
	)
	uncompleteTaskEndpoint := middleware.LoggingMiddleware(log.With(logger, "method", "UncompleteTask"))(
		middleware.MetricsMiddleware(requestLatencyMetric)(
			MakeUncompleteTaskEndpoint(svc),
		),
	)
	deleteTaskEndpoint := middleware.LoggingMiddleware(log.With(logger, "method", "DeleteTask"))(
		middleware.MetricsMiddleware(requestLatencyMetric)(
			MakeDeleteTaskEndpoint(svc),
		),
	)
	return Endpoints{
		HelloEndpoint:          helloEndpoint,
		ListTasksEndpoint:      listTasksEndpoint,
		GetTaskEndpoint:        getTaskEndpoint,
		CreateTaskEndpoint:     createTaskEndpoint,
		UpdateTaskEndpoint:     updateTaskEndpoint,
		CompleteTaskEndpoint:   completeTaskEndpoint,
		UncompleteTaskEndpoint: uncompleteTaskEndpoint,
		DeleteTaskEndpoint:     deleteTaskEndpoint,
	}
}

func MakeHelloEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(HelloRequest)
		greeting, err := svc.Hello(req.Name)
		if err != nil {
			return nil, err
		}
		return HelloResponse{Greeting: greeting}, nil
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
		req := request.(models.Task)
		task, err := svc.CreateTask(&req)
		if err != nil {
			return nil, err
		}
		return *task, nil
	}
}

func MakeUpdateTaskEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(models.Task)
		task, err := svc.UpdateTask(&req)
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
		return struct{}{}, nil
	}
}

func MakeUncompleteTaskEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(UncompleteTaskRequest)
		if err := svc.UncompleteTask(req.ID); err != nil {
			return nil, err
		}
		return struct{}{}, nil
	}
}

func MakeDeleteTaskEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(DeleteTaskRequest)
		if err := svc.DeleteTask(req.ID); err != nil {
			return nil, err
		}
		return struct{}{}, nil
	}
}
