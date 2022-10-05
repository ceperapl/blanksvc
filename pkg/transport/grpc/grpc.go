package grpc

import (
	"github.com/company/blanksvc/pkg/endpoints"

	"github.com/go-kit/kit/transport"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/go-kit/log"
	context "golang.org/x/net/context"
)

type grpcServer struct {
	hello grpctransport.Handler
}

func NewGRPCServer(endpoints endpoints.Endpoints, logger log.Logger) HelloServer {
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

func (g *grpcServer) Hello(ctx context.Context, req *HelloRequest) (*HelloResponse, error) {
	_, resp, err := g.hello.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*HelloResponse), nil
}

func decodeGRPCHelloRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*HelloRequest)
	return endpoints.HelloRequest{Name: req.Name}, nil
}

func encodeGRPCHelloResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(endpoints.HelloResponse)
	return &HelloResponse{Greeting: resp.Greeting}, nil
}
