package main

import (
	"github.com/PatrickKvartsaniy/print-me-at/config"
	"github.com/PatrickKvartsaniy/print-me-at/repository"
	"github.com/PatrickKvartsaniy/print-me-at/server"
)

func main() {
	cfg := config.ReadOS()
	repo := repository.NewRedisRepository(cfg.Redis)
	srv := server.CreateAndRun(cfg.Port)
}
