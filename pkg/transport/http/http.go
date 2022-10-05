package http

import (
	"context"
	"encoding/json"

	"net/http"

	"github.com/company/blanksvc/pkg/endpoints"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
)

func NewHTTPHandler(endpoints endpoints.Endpoints, logger log.Logger) http.Handler {
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(errorEncoder),
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}

	return httptransport.NewServer(
		endpoints.HelloEndpoint,
		decodeHTTPHelloRequest,
		encodeHTTPGenericResponse,
		options...,
	)
}

func decodeHTTPHelloRequest(_ context.Context, r *http.Request) (interface{}, error) {
	name := r.URL.Query().Get("name")
	return endpoints.HelloRequest{Name: name}, nil
}

func encodeHTTPGenericResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if f, ok := response.(endpoint.Failer); ok && f.Failed() != nil {
		errorEncoder(ctx, f.Failed(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
