package middleware

import (
	"context"
	"time"

	"github.com/company/blanksvc/pkg/metrics"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/prometheus/client_golang/prometheus"
)

func MetricsMiddleware(latency *prometheus.HistogramVec) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			defer func(begin time.Time) {
				latency.With(prometheus.Labels{
					metrics.MetricsEndpointLabel: ctx.Value(httptransport.ContextKeyRequestPath).(string),
					metrics.MetricsMethodLabel:   ctx.Value(httptransport.ContextKeyRequestMethod).(string),
				}).Observe(time.Since(begin).Seconds())
			}(time.Now())
			return next(ctx, request)
		}
	}
}
