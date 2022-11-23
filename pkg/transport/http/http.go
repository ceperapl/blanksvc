package http

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"net/http"

	"github.com/company/blanksvc/pkg/endpoints"
	"github.com/company/blanksvc/pkg/models"
	"github.com/gorilla/mux"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
)

func Handle(mux *mux.Router, endpoints endpoints.Endpoints, logger log.Logger) *mux.Router {
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(errorEncoder),
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		httptransport.ServerBefore(httptransport.PopulateRequestContext),
	}

	subRouter := mux.PathPrefix("/api/v1").Subrouter()

	subRouter.Handle("/tasks", httptransport.NewServer(
		endpoints.ListTasksEndpoint,
		decodeHTTPListTasksRequest,
		encodeHTTPGenericResponse,
		options...,
	)).Methods(http.MethodGet)

	subRouter.Handle("/tasks/{id}", httptransport.NewServer(
		endpoints.GetTaskEndpoint,
		decodeHTTPGetTaskRequest,
		encodeHTTPGenericResponse,
		options...,
	)).Methods(http.MethodGet)

	subRouter.Handle("/tasks", httptransport.NewServer(
		endpoints.CreateTaskEndpoint,
		decodeHTTPCreateTaskRequest,
		encodeHTTPGenericResponse,
		options...,
	)).Methods(http.MethodPost)

	subRouter.Handle("/tasks", httptransport.NewServer(
		endpoints.UpdateTaskEndpoint,
		decodeHTTPUpdateTaskRequest,
		encodeHTTPGenericResponse,
		options...,
	)).Methods(http.MethodPut)

	subRouter.Handle("/tasks/{id}/complete", httptransport.NewServer(
		endpoints.CompleteTaskEndpoint,
		decodeHTTPCompleteTaskRequest,
		encodeHTTPGenericResponse,
		options...,
	)).Methods(http.MethodPost)

	subRouter.Handle("/tasks/{id}/uncomplete", httptransport.NewServer(
		endpoints.UncompleteTaskEndpoint,
		decodeHTTPUncompleteTaskRequest,
		encodeHTTPGenericResponse,
		options...,
	)).Methods(http.MethodPost)

	subRouter.Handle("/tasks/{id}", httptransport.NewServer(
		endpoints.DeleteTaskEndpoint,
		decodeHTTPDeleteTaskRequest,
		encodeHTTPGenericResponse,
		options...,
	)).Methods(http.MethodDelete)

	return subRouter
}

func decodeHTTPListTasksRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var err error
	filter := r.URL.Query().Get("filter")
	sort := r.URL.Query().Get("sort")
	var itemsOnPage int
	if _, ok := r.URL.Query()["itemsOnPage"]; ok {
		if itemsOnPage, err = strconv.Atoi(r.URL.Query().Get("itemsOnPage")); err != nil {
			return nil, fmt.Errorf("conversion of itemsOnPage to int: %w", err)
		}
	}
	var page int
	if _, ok := r.URL.Query()["page"]; ok {
		if page, err = strconv.Atoi(r.URL.Query().Get("page")); err != nil {
			return nil, fmt.Errorf("conversion of page to int: %w", err)
		}
	}
	return endpoints.ListTasksRequest{
		Filter:      filter,
		Sort:        sort,
		ItemsOnPage: itemsOnPage,
		Page:        page,
	}, nil
}

func decodeHTTPGetTaskRequest(_ context.Context, r *http.Request) (interface{}, error) {
	id := mux.Vars(r)["id"]
	return endpoints.GetTaskRequest{ID: id}, nil
}

func decodeHTTPCreateTaskRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var task *models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		return nil, fmt.Errorf("task unmarshaling: %w", err)
	}
	return task, nil
}

func decodeHTTPUpdateTaskRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var task *models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		return nil, fmt.Errorf("task unmarshaling: %w", err)
	}
	return task, nil
}

func decodeHTTPCompleteTaskRequest(_ context.Context, r *http.Request) (interface{}, error) {
	id := mux.Vars(r)["id"]
	return endpoints.CompleteTaskRequest{ID: id}, nil
}

func decodeHTTPUncompleteTaskRequest(_ context.Context, r *http.Request) (interface{}, error) {
	id := mux.Vars(r)["id"]
	return endpoints.UncompleteTaskRequest{ID: id}, nil
}

func decodeHTTPDeleteTaskRequest(_ context.Context, r *http.Request) (interface{}, error) {
	id := mux.Vars(r)["id"]
	return endpoints.DeleteTaskRequest{ID: id}, nil
}

func encodeHTTPGenericResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if f, ok := response.(endpoint.Failer); ok && f.Failed() != nil {
		errorEncoder(ctx, f.Failed(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// nolint: wrapcheck
	return json.NewEncoder(w).Encode(response)
}
