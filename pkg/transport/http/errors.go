package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/company/blanksvc/pkg/metrics"
	"github.com/company/blanksvc/pkg/service"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/prometheus/client_golang/prometheus"
)

type errorWrapper struct {
	Error string `json:"error"`
}

func errorEncoder(errorCounterMetric *prometheus.CounterVec) httptransport.ErrorEncoder {
	return func(_ context.Context, err error, w http.ResponseWriter) {
		responseCode := errorToCode(err)

		errorCounterMetric.With(prometheus.Labels{
			metrics.MetricsErrorTypeLabel:    http.StatusText(responseCode),
			metrics.MetricsErrorKindLabel:    errorKind(err),
			metrics.MetricsResponseCodeLabel: strconv.Itoa(responseCode),
		}).Inc()

		w.WriteHeader(responseCode)
		// nolint: errcheck
		json.NewEncoder(w).Encode(errorWrapper{Error: err.Error()})
	}
}

func errorToCode(err error) int {
	switch {
	case errors.Is(err, service.ErrEmptyString):
		return http.StatusBadRequest
	}
	return http.StatusInternalServerError
}

func errorKind(err error) string {
	switch {
	case errors.Is(err, service.ErrEmptyString):
		return "ErrEmptyString"
	default:
		return "ErrUndefined"
	}
}
