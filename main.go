package main

import (
	"context"
	"github.com/PatrickKvartsaniy/print-me-at/config"
	"github.com/PatrickKvartsaniy/print-me-at/health"
	"github.com/PatrickKvartsaniy/print-me-at/repository"
	"github.com/PatrickKvartsaniy/print-me-at/server"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	cfg := config.ReadOS()
	initLogger(cfg.LogLevel, cfg.PrettyLogOutput)

	ctx, cancel := context.WithCancel(context.Background())
	setupGracefulShutdown(cancel)

	repo := repository.NewRedisRepository(cfg.Redis)
	defer repo.Close()

	srv := server.CreateAndRun(ctx, cfg.Port, repo)
	defer closeWithTimeout(srv.Close, 5)

	healthCheckSrv := health.CreateAndRun(cfg.HealthCheckPort, []health.Check{
		repo.HealthCheck,
		srv.HealthCheck,
	})
	defer closeWithTimeout(healthCheckSrv.Close, 5)

	<-ctx.Done()
}

func closeWithTimeout(close func(context.Context), d time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), d)
	defer cancel()
	close(ctx)
}

func setupGracefulShutdown(stop func()) {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChannel
		logrus.Println("Got Interrupt signal")
		stop()
	}()
}

func initLogger(logLevel string, pretty bool) {
	if pretty {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}
	logrus.SetOutput(os.Stderr)

	switch strings.ToLower(logLevel) {
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	default:
		logrus.SetLevel(logrus.DebugLevel)
	}
}
