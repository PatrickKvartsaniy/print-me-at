package main

import (
	"context"
	"github.com/PatrickKvartsaniy/print-me-at/config"
	"github.com/PatrickKvartsaniy/print-me-at/health"
	"github.com/PatrickKvartsaniy/print-me-at/repository"
	"github.com/PatrickKvartsaniy/print-me-at/server"
	"sync"
)

func main() {
	ctx := context.Background()
	cfg := config.ReadOS()

	repo := repository.NewRedisRepository(cfg.Redis)
	defer repo.Close()

	srv := server.CreateAndRun(ctx, cfg.Port, repo)
	defer srv.Close(ctx)

	healthCheckSrv := health.New(cfg.HealthCheckPort, []health.Check{
		repo.HealthCheck,
		srv.HealthCheck,
	})
	var wg sync.WaitGroup
	healthCheckSrv.Run(ctx, &wg)
	wg.Wait()
}
