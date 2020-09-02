package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/PatrickKvartsaniy/print-me-at/models"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Repository interface {
	AddNewTask(msg models.Message) error
	Poll(ctx context.Context)
}

type Server struct {
	http      *http.Server
	runErr    error
	readiness bool
	repo      Repository
}

func CreateAndRun(ctx context.Context, port int, repo Repository) *Server {
	service := &Server{
		http: &http.Server{
			Addr: fmt.Sprintf(":%d", port),
		},
		repo: repo,
	}

	service.setupHandlers()

	service.run()
	service.repo.Poll(ctx)

	return service
}

func (s *Server) setupHandlers() {
	h := http.NewServeMux()
	handler := http.HandlerFunc(s.messageScheduler)
	h.Handle("/printMeAt", loggingMiddleware(handler))
	s.http.Handler = h
}

func (s *Server) run() {
	log.Info("server is running")

	go func() {
		log.Debug("http service: addr=", s.http.Addr)
		if err := s.http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.runErr = err
			log.WithError(err).Error("http service")
		}
	}()
	s.readiness = true
}

func (s *Server) Close(ctx context.Context) {
	if err := s.http.Shutdown(ctx); err != nil {
		log.WithError(err).Error("stopping http service")
	}
	log.Info("http service stopped")
}

func (s *Server) HealthCheck() error {
	if !s.readiness {
		return errors.New("http service isn't ready yet")
	}
	if s.runErr != nil {
		return errors.New("http service: run issue")
	}
	return nil
}
