package cmd

import (
	"fmt"
	"net"
	"net/http"
	"os"

	hellov1 "github.com/company/blanksvc/gen/proto/go/hello/v1"
	"github.com/company/blanksvc/pkg/endpoints"
	"github.com/company/blanksvc/pkg/service"
	grpctransport "github.com/company/blanksvc/pkg/transport/grpc"
	httptransport "github.com/company/blanksvc/pkg/transport/http"
	"github.com/company/blanksvc/pkg/utils/healthcheck"
	"github.com/prometheus/client_golang/prometheus/promhttp"

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
	// Configure metrics
	rootMux.Handle("/metrics", promhttp.Handler())
	// Configure REST API
	subrouter := rootMux.PathPrefix("/api/v1").Subrouter()
	subrouter.Handle("/hello", httpHandler)
	// Start the HTTP server
	httpServerAddr := fmt.Sprintf("0.0.0.0:%d", config.HTTPServer.Port)
	// nolint: errcheck
	logger.Log("transport", "HTTP", "addr", httpServerAddr)
	go func() {
		doneC <- http.ListenAndServe(httpServerAddr, rootMux)
	}()

	// Configure the GRPC server
	grpcServerAddr := fmt.Sprintf("0.0.0.0:%d", config.GRPCServer.Port)
	grpcListener, err := net.Listen("tcp", grpcServerAddr)
	if err != nil {
		// nolint: errcheck
		logger.Log("transport", "gRPC", "during", "Listen", "err", err)
		return err
	}
	baseServer := grpc.NewServer()
	hellov1.RegisterHelloServiceServer(baseServer, grpcServer)
	// Start the GRPC server
	// nolint: errcheck
	logger.Log("transport", "GRPC", "addr", grpcServerAddr)
	go func() {
		doneC <- baseServer.Serve(grpcListener)
	}()

	// waiting for the errors from servers
	if err := <-doneC; err != nil {
		// nolint: errcheck
		logger.Log("err", err)
		return err
	}

	return nil
}

func readinessCheck() error {
	// return errors.New("error")
	return nil
}
