package middleware

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/company/blanksvc/pkg/common"
	"github.com/company/blanksvc/pkg/metrics"
	"github.com/go-kit/kit/endpoint"
	"github.com/prometheus/client_golang/prometheus"
)

func MetricsMiddleware(latency prometheus.Observer, requestCounter prometheus.Counter, errorCounter *prometheus.CounterVec) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			defer func(begin time.Time) {
				latency.Observe(time.Since(begin).Seconds())
				requestCounter.Inc()
			}(time.Now())
			resp, err := next(ctx, request)
			if err != nil {
				responseCode := common.ErrorToCode(err)
				errorCounter.With(prometheus.Labels{
					metrics.MetricsErrorTypeLabel:    http.StatusText(responseCode),
					metrics.MetricsErrorKindLabel:    common.ErrorKind(err),
					metrics.MetricsResponseCodeLabel: strconv.Itoa(responseCode),
				}).Inc()
			}
			return resp, err
		}
	}
}
