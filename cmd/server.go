package cmd

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/company/blanksvc/pkg/endpoints"
	"github.com/company/blanksvc/pkg/service"
	grpctransport "github.com/company/blanksvc/pkg/transport/grpc"
	httptransport "github.com/company/blanksvc/pkg/transport/http"
	"github.com/company/blanksvc/pkg/utils/healthcheck"

	"github.com/go-kit/log"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

func RunServer() error {

	doneC := make(chan error)

	// Init config
	config := NewConfig()

	// Create a single logger, which we'll use and give to other components.
	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)

	// Build the layers of the service "onion" from the inside out
	service := service.New(logger)
	endpoints := endpoints.New(service, logger)
	httpHandler := httptransport.NewHTTPHandler(endpoints, logger)
	grpcServer := grpctransport.NewGRPCServer(endpoints, logger)

	// Configure the HTTP server
	rootMux := mux.NewRouter()
	// Configure health checks
	healthchecker := healthcheck.New()
	healthchecker.AddReadinessChecks(readinessCheck)
	rootMux.Handle(config.HTTPServer.ReadinessEndpoint, healthchecker.ReadinessHandler())
	rootMux.Handle(config.HTTPServer.LivenessEndpoint, healthchecker.LivenessHandler())
	// Configure REST API
	subrouter := rootMux.PathPrefix("/api/v1").Subrouter()
	subrouter.Handle("/hello", httpHandler)
	// Start the HTTP server
	httpServerAddr := fmt.Sprintf("0.0.0.0:%d", config.HTTPServer.Port)
	logger.Log("transport", "HTTP", "addr", httpServerAddr)
	go func() {
		doneC <- http.ListenAndServe(httpServerAddr, rootMux)
	}()

	// Configure the GRPC server
	grpcServerAddr := fmt.Sprintf("0.0.0.0:%d", config.GRPCServer.Port)
	grpcListener, err := net.Listen("tcp", grpcServerAddr)
	if err != nil {
		logger.Log("transport", "gRPC", "during", "Listen", "err", err)
		return err
	}
	baseServer := grpc.NewServer()
	grpctransport.RegisterHelloServer(baseServer, grpcServer)
	// Start the GRPC server
	logger.Log("transport", "GRPC", "addr", grpcServerAddr)
	go func() {
		doneC <- baseServer.Serve(grpcListener)
	}()

	// waiting for the errors from servers
	if err := <-doneC; err != nil {
		logger.Log("err", err)
		return err
	}

	return nil
}

func readinessCheck() error {
	// return errors.New("error")
	return nil
}
