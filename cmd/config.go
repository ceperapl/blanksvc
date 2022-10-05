package cmd

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	envPrefix = "BLANK"

	httpServerPortEnv     = "HTTP_SERVER_PORT"
	httpServerPortDefault = 8080

	httpReadinessEndpointEnv     = "HTTP_READINESS_ENDPOINT"
	httpReadinessEndpointDefault = "/ready"

	httpLivenessEndpointEnv     = "HTTP_LIVENESS_ENDPOINT"
	httpLivenessEndpointDefault = "/health"

	grpcServerPortEnv     = "GRPC_SERVER_PORT"
	grpcServerPortDefault = 9000
)

type Config struct {
	HTTPServer struct {
		Port              uint16
		ReadinessEndpoint string
		LivenessEndpoint  string
	}
	GRPCServer struct {
		Port uint16
	}
}

func NewConfig() Config {
	config := Config{}

	// Init config via env variables
	viper.SetEnvPrefix("PRF")

	viper.BindEnv(httpServerPortEnv)
	viper.SetDefault(httpServerPortEnv, httpServerPortDefault)
	config.HTTPServer.Port = viper.GetUint16(httpServerPortEnv)

	viper.BindEnv(httpReadinessEndpointEnv)
	viper.SetDefault(httpReadinessEndpointEnv, httpReadinessEndpointDefault)
	config.HTTPServer.ReadinessEndpoint = viper.GetString(httpReadinessEndpointEnv)

	viper.BindEnv(httpLivenessEndpointEnv)
	viper.SetDefault(httpLivenessEndpointEnv, httpLivenessEndpointDefault)
	config.HTTPServer.LivenessEndpoint = viper.GetString(httpLivenessEndpointEnv)

	viper.BindEnv(grpcServerPortEnv)
	viper.SetDefault(grpcServerPortEnv, grpcServerPortDefault)
	config.GRPCServer.Port = viper.GetUint16(grpcServerPortEnv)

	// Init config via flags
	pflag.Uint16Var(&config.HTTPServer.Port, "httpserver.port", config.HTTPServer.Port, fmt.Sprintf("HTTP Server port; env: %s", httpServerPortEnv))
	pflag.StringVar(&config.HTTPServer.ReadinessEndpoint, "httpserver.readiness", config.HTTPServer.ReadinessEndpoint, fmt.Sprintf("HTTP Server readiness endpoint name; env: %s", httpReadinessEndpointEnv))
	pflag.StringVar(&config.HTTPServer.LivenessEndpoint, "httpserver.liveness", config.HTTPServer.LivenessEndpoint, fmt.Sprintf("HTTP Server liveness endpoint name; env: %s", httpLivenessEndpointEnv))
	pflag.Uint16Var(&config.GRPCServer.Port, "grpcserver.port", config.GRPCServer.Port, fmt.Sprintf("gRPC Server port; env: %s", grpcServerPortEnv))

	pflag.Parse()

	return config
}
