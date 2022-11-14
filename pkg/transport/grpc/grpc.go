package grpc

import (
	"context"
	"time"

	"github.com/company/blanksvc/pkg/endpoints"
	"github.com/company/blanksvc/pkg/models"
	"google.golang.org/protobuf/types/known/timestamppb"

	taskv1 "github.com/company/blanksvc/gen/proto/go/task/v1"
	"github.com/go-kit/kit/transport"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/go-kit/log"
)

type ServiceServer interface {
	taskv1.TaskServiceServer
}

type grpcServer struct {
	listTasks      grpctransport.Handler
	getTask        grpctransport.Handler
	createTask     grpctransport.Handler
	updateTask     grpctransport.Handler
	completeTask   grpctransport.Handler
	uncompleteTask grpctransport.Handler
	deleteTask     grpctransport.Handler
}

func NewGRPCServer(endpoints endpoints.Endpoints, logger log.Logger) ServiceServer {
	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}

	return &grpcServer{
		listTasks: grpctransport.NewServer(
			endpoints.ListTasksEndpoint,
			decodeGRPCListTasksRequest,
			encodeGRPCListTasksResponse,
			options...,
		),
		getTask: grpctransport.NewServer(
			endpoints.GetTaskEndpoint,
			decodeGRPCGetTaskRequest,
			encodeGRPCGetTaskResponse,
			options...,
		),
		createTask: grpctransport.NewServer(
			endpoints.CreateTaskEndpoint,
			decodeGRPCCreateTaskRequest,
			encodeGRPCCreateTaskResponse,
			options...,
		),
		updateTask: grpctransport.NewServer(
			endpoints.UpdateTaskEndpoint,
			decodeGRPCUpdateTaskRequest,
			encodeGRPCUpdateTaskResponse,
			options...,
		),
		completeTask: grpctransport.NewServer(
			endpoints.CompleteTaskEndpoint,
			decodeGRPCCompleteTaskRequest,
			encodeGRPCCompleteTaskResponse,
			options...,
		),
		uncompleteTask: grpctransport.NewServer(
			endpoints.UncompleteTaskEndpoint,
			decodeGRPCUncompleteTaskRequest,
			encodeGRPCUncompleteTaskResponse,
			options...,
		),
		deleteTask: grpctransport.NewServer(
			endpoints.DeleteTaskEndpoint,
			decodeGRPCDeleteTaskRequest,
			encodeGRPCDeleteTaskResponse,
			options...,
		),
	}
}

func (g *grpcServer) ListTasks(ctx context.Context, req *taskv1.ListTasksRequest) (*taskv1.ListTasksResponse, error) {
	_, resp, err := g.listTasks.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*taskv1.ListTasksResponse), nil
}

func (g *grpcServer) GetTask(ctx context.Context, req *taskv1.GetTaskRequest) (*taskv1.GetTaskResponse, error) {
	_, resp, err := g.getTask.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*taskv1.GetTaskResponse), nil
}

func (g *grpcServer) CreateTask(ctx context.Context, req *taskv1.CreateTaskRequest) (*taskv1.CreateTaskResponse, error) {
	_, resp, err := g.createTask.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*taskv1.CreateTaskResponse), nil
}

func (g *grpcServer) UpdateTask(ctx context.Context, req *taskv1.UpdateTaskRequest) (*taskv1.UpdateTaskResponse, error) {
	_, resp, err := g.updateTask.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*taskv1.UpdateTaskResponse), nil
}

func (g *grpcServer) CompleteTask(ctx context.Context, req *taskv1.CompleteTaskRequest) (*taskv1.CompleteTaskResponse, error) {
	_, resp, err := g.completeTask.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*taskv1.CompleteTaskResponse), nil
}

func (g *grpcServer) UncompleteTask(ctx context.Context, req *taskv1.UncompleteTaskRequest) (*taskv1.UncompleteTaskResponse, error) {
	_, resp, err := g.uncompleteTask.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*taskv1.UncompleteTaskResponse), nil
}

func (g *grpcServer) DeleteTask(ctx context.Context, req *taskv1.DeleteTaskRequest) (*taskv1.DeleteTaskResponse, error) {
	_, resp, err := g.deleteTask.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*taskv1.DeleteTaskResponse), nil
}

func decodeGRPCListTasksRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*taskv1.ListTasksRequest)
	return endpoints.ListTasksRequest{
		Filter:      req.Filter,
		Sort:        req.Sort,
		ItemsOnPage: int(req.ItemsOnPage),
		Page:        int(req.Page),
	}, nil
}

func encodeGRPCListTasksResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(endpoints.ListTaskResponse)
	var grpcTasks []*taskv1.Task
	for _, task := range resp.Tasks {
		task := task
		grpcTasks = append(grpcTasks, TransformTask(&task))
	}

	return &taskv1.ListTasksResponse{
		Tasks: grpcTasks,
		Count: resp.Count,
	}, nil
}

func decodeGRPCGetTaskRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*taskv1.GetTaskRequest)
	return endpoints.GetTaskRequest{
		ID: req.Id,
	}, nil
}

func encodeGRPCGetTaskResponse(_ context.Context, response interface{}) (interface{}, error) {
	task := response.(models.Task)
	return &taskv1.GetTaskResponse{Task: TransformTask(&task)}, nil
}

func decodeGRPCCreateTaskRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*taskv1.CreateTaskRequest)
	return TransformTaskPB(req.Task), nil
}

func encodeGRPCCreateTaskResponse(_ context.Context, response interface{}) (interface{}, error) {
	task := response.(models.Task)
	return &taskv1.CreateTaskResponse{Task: TransformTask(&task)}, nil
}

func decodeGRPCUpdateTaskRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*taskv1.UpdateTaskRequest)
	return TransformTaskPB(req.Task), nil
}

func encodeGRPCUpdateTaskResponse(_ context.Context, response interface{}) (interface{}, error) {
	task := response.(models.Task)
	return &taskv1.UpdateTaskResponse{Task: TransformTask(&task)}, nil
}

func decodeGRPCCompleteTaskRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*taskv1.CompleteTaskRequest)
	return endpoints.CompleteTaskRequest{ID: req.Id}, nil
}

func encodeGRPCCompleteTaskResponse(_ context.Context, response interface{}) (interface{}, error) {
	return &taskv1.CompleteTaskResponse{}, nil
}

func decodeGRPCUncompleteTaskRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*taskv1.UncompleteTaskRequest)
	return endpoints.UncompleteTaskRequest{ID: req.Id}, nil
}

func encodeGRPCUncompleteTaskResponse(_ context.Context, response interface{}) (interface{}, error) {
	return &taskv1.UncompleteTaskResponse{}, nil
}

func decodeGRPCDeleteTaskRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*taskv1.DeleteTaskRequest)
	return endpoints.DeleteTaskRequest{ID: req.Id}, nil
}

func encodeGRPCDeleteTaskResponse(_ context.Context, response interface{}) (interface{}, error) {
	return &taskv1.DeleteTaskResponse{}, nil
}

func TransformTask(task *models.Task) *taskv1.Task {

	var description string
	if task.Description != nil {
		description = *task.Description
	}

	var completedAt *timestamppb.Timestamp
	if task.CompletedAt != nil {
		completedAt = timestamppb.New(*task.CompletedAt)
	}

	var updatedAt *timestamppb.Timestamp
	if task.UpdatedAt != nil {
		updatedAt = timestamppb.New(*task.UpdatedAt)
	}

	taskPB := taskv1.Task{
		Id:          task.ID,
		Name:        task.Name,
		Description: description,
		Deadline:    task.Deadline,
		CompletedAt: completedAt,
		CreatedAt:   timestamppb.New(task.CreatedAt.Local()),
		UpdatedAt:   updatedAt,
	}

	return &taskPB
}

func TransformTaskPB(taskPB *taskv1.Task) *models.Task {

	var description *string
	if taskPB.Description != "" {
		description = &taskPB.Description
	}

	var completedAt *time.Time
	if taskPB.CompletedAt != nil {
		temp := taskPB.CompletedAt.AsTime()
		completedAt = &temp
	}

	var updatedAt *time.Time
	if taskPB.UpdatedAt != nil {
		temp := taskPB.UpdatedAt.AsTime()
		updatedAt = &temp
	}

	task := models.Task{
		ID:          taskPB.Id,
		Name:        taskPB.Name,
		Description: description,
		Deadline:    taskPB.Deadline,
		CompletedAt: completedAt,
		CreatedAt:   taskPB.CreatedAt.AsTime(),
		UpdatedAt:   updatedAt,
	}

	return &task
}
