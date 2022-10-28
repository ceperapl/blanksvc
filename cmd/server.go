package cmd

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	hellov1 "github.com/company/blanksvc/gen/proto/go/hello/v1"
	"github.com/company/blanksvc/pkg/endpoints"
	"github.com/company/blanksvc/pkg/metrics"
	"github.com/company/blanksvc/pkg/repository/postgres"
	"github.com/company/blanksvc/pkg/service"
	grpctransport "github.com/company/blanksvc/pkg/transport/grpc"
	httptransport "github.com/company/blanksvc/pkg/transport/http"
	"github.com/company/blanksvc/pkg/utils"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/go-kit/log"
	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/v4/stdlib"
	"google.golang.org/grpc"
)

func RunServer() error {

	doneC := make(chan error)

	// Init config
	config := NewConfig()

	// Create metrics
	requestLatencyMetric := metrics.NewRequestLatencyMetric()
	requestCounterMetric := metrics.NewRequestCounterMetric()
	errorCounterMetric := metrics.NewErrorCounterMetric()

	// Create a single logger, which we'll use and give to other components.
	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)

	// Connect to Postgres with retry
	var db *sql.DB
	if err := utils.Retry(time.Second, 30, func() (stop bool) {
		logger.Log("database", config.Postgres.DSN)
		var errOpen error
		if db, errOpen = sql.Open("postgres", config.Postgres.DSN); errOpen != nil {
			return false
		}
		if err := db.Ping(); err != nil {
			return false
		}
		return true
	}); err != nil {
		return err
	}
	defer db.Close()

	// Migrate DB
	if err := postgres.Migrate(config.Postgres.DSN); err != nil {
		return err
	}

	// Configure the HTTP server
	rootMux := mux.NewRouter()

	// Build the layers of the service "onion" from the inside out
	repository := postgres.NewRepo()
	service := service.New(logger, repository)
	endpoints := endpoints.New(service, logger, requestLatencyMetric)
	httptransport.Handle(rootMux, endpoints, logger, requestCounterMetric, errorCounterMetric)
	grpcServer := grpctransport.NewGRPCServer(endpoints, logger)

	// Configure health checks
	healthchecker := utils.New()
	healthchecker.AddReadinessChecks(readinessCheck)
	rootMux.Handle(config.HTTPServer.ReadinessEndpoint, healthchecker.ReadinessHandler())
	rootMux.Handle(config.HTTPServer.LivenessEndpoint, healthchecker.LivenessHandler())
	// Configure metrics
	rootMux.Handle("/metrics", promhttp.Handler())
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
