package metrics

import (
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/prometheus"
	kitprometheus "github.com/prometheus/client_golang/prometheus"
)

var (
	MetricsNamespace = "blanksvc"

	MetricsEndpointNameLabel = "endpoint_name"
	MetricsResponseCodeLabel = "response_code"
	MetricsErrorTypeLabel    = "error_type"
	MetricsErrorKindLabel    = "error_kind"
)

type Metrics struct {
	RequestLatency metrics.Histogram
	RequestCounter metrics.Counter
	ErrorCounter   metrics.Counter
}

func NewRequestLatencyMetric() metrics.Histogram {
	latencyHistogram := prometheus.NewHistogramFrom(kitprometheus.HistogramOpts{
		Namespace: MetricsNamespace,
		Subsystem: "service",
		Name:      "request_latency",
		Help:      "Latency of request processing",
	}, []string{MetricsEndpointNameLabel})
	return latencyHistogram
}

func NewRequestCounterMetric() metrics.Counter {
	requestCounter := prometheus.NewCounterFrom(kitprometheus.CounterOpts{
		Namespace: MetricsNamespace,
		Subsystem: "service",
		Name:      "requests_total",
		Help:      "Amount of requests",
	}, []string{MetricsEndpointNameLabel})
	return requestCounter
}

func NewErrorCounterMetric() metrics.Counter {
	errorCounter := prometheus.NewCounterFrom(kitprometheus.CounterOpts{
		Namespace: MetricsNamespace,
		Subsystem: "service",
		Name:      "errors_total",
		Help:      "Amount of errors",
	}, []string{MetricsErrorTypeLabel, MetricsErrorKindLabel, MetricsResponseCodeLabel})
	return errorCounter
}
