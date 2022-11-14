package metrics

import (
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	MetricsNamespace = "blanksvc"

	MetricsEndpointNameLabel = "endpoint_name"
	MetricsResponseCodeLabel = "response_code"
	MetricsErrorTypeLabel    = "error_type"
	MetricsErrorKindLabel    = "error_kind"
)

// nolint: gochecknoglobals
var (
	// RequestLatency is a metric for a request processing latency
	RequestLatency = kitprometheus.NewHistogramFrom(prometheus.HistogramOpts{
		Namespace: MetricsNamespace,
		Subsystem: "service",
		Name:      "request_latency",
		Help:      "Latency of request processing",
	}, []string{MetricsEndpointNameLabel})

	// RequestCounter is a metric for amount of requests
	RequestCounter = kitprometheus.NewCounterFrom(prometheus.CounterOpts{
		Namespace: MetricsNamespace,
		Subsystem: "service",
		Name:      "requests_total",
		Help:      "Amount of requests",
	}, []string{MetricsEndpointNameLabel})

	// ErrorCounter is a metric for amount of errors
	ErrorCounter = kitprometheus.NewCounterFrom(prometheus.CounterOpts{
		Namespace: MetricsNamespace,
		Subsystem: "service",
		Name:      "errors_total",
		Help:      "Amount of errors",
	}, []string{MetricsErrorTypeLabel, MetricsErrorKindLabel, MetricsResponseCodeLabel})
)
