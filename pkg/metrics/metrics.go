package metrics

import (
	prometheus "github.com/prometheus/client_golang/prometheus"
)

var (
	MetricsNamespace = "blanksvc"

	MetricsEndpointNameLabel = "endpoint_name"
	MetricsResponseCodeLabel = "response_code"
	MetricsErrorTypeLabel    = "error_type"
	MetricsErrorKindLabel    = "error_kind"
)

type RequestLatencyMetric struct {
	metric *prometheus.HistogramVec
}

type RequestCounterMetric struct {
	metric *prometheus.CounterVec
}

func NewRequestLatencyMetric() RequestLatencyMetric {
	latencyHistogram := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: MetricsNamespace,
		Subsystem: "service",
		Name:      "request_latency",
		Help:      "Latency of request processing",
	}, []string{MetricsEndpointNameLabel})
	prometheus.DefaultRegisterer.MustRegister(latencyHistogram)
	return RequestLatencyMetric{latencyHistogram}
}

func (r *RequestLatencyMetric) SetLabels(metricsEndpointNameLabel string) prometheus.Observer {
	return r.metric.With(prometheus.Labels{
		MetricsEndpointNameLabel: metricsEndpointNameLabel,
	})
}

func NewRequestCounterMetric() RequestCounterMetric {
	requestCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: MetricsNamespace,
		Subsystem: "service",
		Name:      "requests_total",
		Help:      "Amount of requests",
	}, []string{MetricsEndpointNameLabel})
	prometheus.DefaultRegisterer.MustRegister(requestCounter)
	return RequestCounterMetric{requestCounter}
}

func (r *RequestCounterMetric) SetLabels(metricsEndpointNameLabel string) prometheus.Counter {
	return r.metric.With(prometheus.Labels{
		MetricsEndpointNameLabel: metricsEndpointNameLabel,
	})
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
