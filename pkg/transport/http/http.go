package http

import (
	"context"
	"encoding/json"

	"net/http"

	"github.com/company/blanksvc/pkg/endpoints"
	"github.com/company/blanksvc/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
)

func NewHTTPHandler(endpoints endpoints.Endpoints, logger log.Logger, requestCounterMetric, errorCounterMetric *prometheus.CounterVec) http.Handler {
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(errorEncoder(errorCounterMetric)),
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		httptransport.ServerBefore(httptransport.PopulateRequestContext),
		httptransport.ServerFinalizer(
			metrics.MetricsOption(requestCounterMetric),
		),
	}

	return httptransport.NewServer(
		endpoints.HelloEndpoint,
		decodeHTTPHelloRequest,
		encodeHTTPGenericResponse(errorCounterMetric),
		options...,
	)
}

func decodeHTTPHelloRequest(_ context.Context, r *http.Request) (interface{}, error) {
	name := r.URL.Query().Get("name")
	return endpoints.HelloRequest{Name: name}, nil
}

func encodeHTTPGenericResponse(errorCounterMetric *prometheus.CounterVec) httptransport.EncodeResponseFunc {
	return func(ctx context.Context, w http.ResponseWriter, response interface{}) error {
		if f, ok := response.(endpoint.Failer); ok && f.Failed() != nil {
			errorEncoder(errorCounterMetric)(ctx, f.Failed(), w)
			return nil
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		return json.NewEncoder(w).Encode(response)
	}
}
