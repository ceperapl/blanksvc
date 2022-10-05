package endpoints

import (
	"context"

	"github.com/company/blanksvc/pkg/endpoints/middleware"
	"github.com/company/blanksvc/pkg/service"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
)

type HelloRequest struct {
	Name string `json:"name"`
}

type HelloResponse struct {
	Greeting string `json:"greeting"`
}

type Endpoints struct {
	HelloEndpoint endpoint.Endpoint
}

func New(svc service.Service, logger log.Logger) Endpoints {
	return Endpoints{
		HelloEndpoint: middleware.LoggingMiddleware(log.With(logger, "method", "Hello"))(MakeHelloEndpoint(svc)),
	}
}

func MakeHelloEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(HelloRequest)
		greeting, err := svc.Hello(req.Name)
		if err != nil {
			return nil, err
		}
		return HelloResponse{Greeting: greeting}, nil
	}
}
