package middleware

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/company/blanksvc/pkg/common"
	"github.com/company/blanksvc/pkg/metrics"
	"github.com/go-kit/kit/endpoint"
	kitmetrics "github.com/go-kit/kit/metrics"
)

func MetricsMiddleware(latency kitmetrics.Histogram, requestCounter kitmetrics.Counter, errorCounter kitmetrics.Counter) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			defer func(begin time.Time) {
				latency.Observe(time.Since(begin).Seconds())
				requestCounter.Add(1)
			}(time.Now())
			resp, err := next(ctx, request)
			if err != nil {
				responseCode := common.ErrorToCode(err)
				errorCounter.With(
					metrics.MetricsErrorTypeLabel, http.StatusText(responseCode),
					metrics.MetricsErrorKindLabel, common.ErrorKind(err),
					metrics.MetricsResponseCodeLabel, strconv.Itoa(responseCode),
				).Add(1)
			}
			return resp, err
		}
	}
}
