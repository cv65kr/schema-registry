package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/cv65kr/schema-registry/internal/grpc"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	// Flags
	fs := pflag.NewFlagSet("default", pflag.ContinueOnError)
	fs.Int("grpc-port", 9999, "gRPC port")
	fs.Duration("server-shutdown-timeout", 10*time.Second, "server graceful shutdown timeout duration")
	fs.String("level", "info", "log level debug, info, warn, error, fatal or panic")

	viper.BindPFlags(fs)
	viper.SetEnvPrefix("SCHEMA_REGISTRY")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	// Logger
	logLevel, err := zap.ParseAtomicLevel(viper.GetString("level"))
	if err != nil {
		fmt.Println("Invalid logger level")
		os.Exit(1)
	}

	config := zap.NewProductionConfig()
	config.Level.SetLevel(logLevel.Level())
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, err := config.Build()
	if err != nil {
		fmt.Println("Logger is not initiated")
		os.Exit(1)
	}
	defer logger.Sync()

	var grpcCfg grpc.Config
	if err := viper.Unmarshal(&grpcCfg); err != nil {
		logger.Panic("onfig unmarshal failed", zap.Error(err))
	}

	// Signal handling
	stopCh := grpc.RunSignalHandler()

	// Start gRPC server
	grpcSrv, _ := grpc.NewServer(&grpcCfg, logger)
	grpcServer := grpcSrv.ListenAndServe()

	// Graceful shutdown
	sd, _ := grpc.NewShutdown(grpcCfg.ServerShutdownTimeout, logger)
	sd.GracefulShutdown(stopCh, grpcServer)
}
