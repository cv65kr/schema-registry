package grpc

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Shutdown struct {
	logger                *zap.Logger
	serverShutdownTimeout time.Duration
}

func NewShutdown(serverShutdownTimeout time.Duration, logger *zap.Logger) (*Shutdown, error) {
	srv := &Shutdown{
		logger:                logger,
		serverShutdownTimeout: serverShutdownTimeout,
	}

	return srv, nil
}

func (s *Shutdown) GracefulShutdown(stopCh <-chan struct{}, grpcServer *grpc.Server) {
	ctx := context.Background()

	// Waiting for SIGNALS
	<-stopCh

	_, cancel := context.WithTimeout(ctx, s.serverShutdownTimeout)
	defer cancel()

	s.logger.Info("Shutting down GRPC server", zap.Duration("timeout", s.serverShutdownTimeout))
	grpcServer.GracefulStop()
}
