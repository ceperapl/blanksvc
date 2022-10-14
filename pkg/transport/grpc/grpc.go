package grpc

import (
	"context"

	"github.com/company/blanksvc/pkg/endpoints"

	hellov1 "github.com/company/blanksvc/gen/proto/go/hello/v1"
	"github.com/go-kit/kit/transport"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/go-kit/log"
)

type grpcServer struct {
	hello grpctransport.Handler
}

func NewGRPCServer(endpoints endpoints.Endpoints, logger log.Logger) hellov1.HelloServiceServer {
	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}

	return &grpcServer{
		hello: grpctransport.NewServer(
			endpoints.HelloEndpoint,
			decodeGRPCHelloRequest,
			encodeGRPCHelloResponse,
			options...,
		),
	}
}

func (g *grpcServer) Hello(ctx context.Context, req *hellov1.HelloRequest) (*hellov1.HelloResponse, error) {
	_, resp, err := g.hello.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*hellov1.HelloResponse), nil
}

func decodeGRPCHelloRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*hellov1.HelloRequest)
	return endpoints.HelloRequest{Name: req.Name}, nil
}

func encodeGRPCHelloResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(endpoints.HelloResponse)
	return &hellov1.HelloResponse{Greeting: resp.Greeting}, nil
}
