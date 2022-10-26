package metrics

import (
	"context"
	"net/http"
	"strconv"

	httptransport "github.com/go-kit/kit/transport/http"
	prometheus "github.com/prometheus/client_golang/prometheus"
)

var (
	MetricsNamespace = "blanksvc"

	MetricsEndpointLabel     = "endpoint"
	MetricsMethodLabel       = "method"
	MetricsResponseCodeLabel = "response_code"
	MetricsErrorTypeLabel    = "error_type"
	MetricsErrorKindLabel    = "error_kind"
)

func NewRequestLatencyMetric() *prometheus.HistogramVec {
	latencyHistogram := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: MetricsNamespace,
		Subsystem: "service",
		Name:      "request_latency",
		Help:      "Latency of request processing",
	}, []string{MetricsEndpointLabel, MetricsMethodLabel})
	prometheus.DefaultRegisterer.MustRegister(latencyHistogram)
	return latencyHistogram
}

func NewRequestCounterMetric() *prometheus.CounterVec {
	requestCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: MetricsNamespace,
		Subsystem: "service",
		Name:      "requests_total",
		Help:      "Amount of requests",
	}, []string{MetricsEndpointLabel, MetricsMethodLabel, MetricsResponseCodeLabel})
	prometheus.DefaultRegisterer.MustRegister(requestCounter)
	return requestCounter
}

func NewErrorCounterMetric() *prometheus.CounterVec {
	errorCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: MetricsNamespace,
		Subsystem: "service",
		Name:      "errors_total",
		Help:      "Amount of errors",
	}, []string{MetricsErrorTypeLabel, MetricsErrorKindLabel, MetricsResponseCodeLabel})
	prometheus.DefaultRegisterer.MustRegister(errorCounter)
	return errorCounter
}

func MetricsOption(requestCounterMetric *prometheus.CounterVec) httptransport.ServerFinalizerFunc {
	return func(ctx context.Context, code int, _ *http.Request) {
		requestCounterMetric.With(prometheus.Labels{
			MetricsEndpointLabel:     ctx.Value(httptransport.ContextKeyRequestPath).(string),
			MetricsMethodLabel:       ctx.Value(httptransport.ContextKeyRequestMethod).(string),
			MetricsResponseCodeLabel: strconv.Itoa(code),
		}).Inc()
	}
}
