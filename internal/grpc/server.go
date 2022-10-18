package grpc

import (
	"fmt"
	"net"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type Config struct {
	Port                  int           `mapstructure:"grpc-port"`
	ServerShutdownTimeout time.Duration `mapstructure:"server-shutdown-timeout"`
}

type Server struct {
	logger *zap.Logger
	config *Config
}

func NewServer(config *Config, logger *zap.Logger) (*Server, error) {
	srv := &Server{
		config: config,
		logger: logger,
	}

	return srv, nil
}

func (s *Server) ListenAndServe() *grpc.Server {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", s.config.Port))
	if err != nil {
		s.logger.Fatal("Failed to listen", zap.Int("port", s.config.Port))
	}

	srv := grpc.NewServer()
	server := health.NewServer()
	reflection.Register(srv)
	grpc_health_v1.RegisterHealthServer(srv, server)
	server.SetServingStatus("Schema registry", grpc_health_v1.HealthCheckResponse_SERVING)

	go func() {
		if err := srv.Serve(listener); err != nil {
			s.logger.Fatal("Failed to serve", zap.Error(err))
		}
	}()

	s.logger.Info("gRPC server started on port", zap.Int("port", s.config.Port))

	return srv
}
