package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type Repository interface {

}

type Server struct {
	http      *http.Server
	runErr    error
	readiness bool
	repo Repository
}

func CreateAndRun(port int) *Server {
	service := &Server{
		http: &http.Server{
			Addr: fmt.Sprintf(":%d", port),
		},
	}

	service.setupHandlers()

	service.run()

	return service
}

func (s *Server) setupHandlers() {
	h := http.NewServeMux()
	handler := http.HandlerFunc(messageScheduler)
	h.Handle("/printMeAt", correlationID(loggingMiddleware(handler)))
	s.http.Handler = h
}

func (s *Server) run() {
	log.Info("server is running")

	go func() {
		log.Debug("graphql service: addr=", s.http.Addr)
		if err := s.http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.runErr = err
			log.WithError(err).Error("graphql service")
		}
	}()
	s.readiness = true
}

func (s *Server) Close(ctx context.Context) {
	if err := s.http.Shutdown(ctx); err != nil {
		log.WithError(err).Error("stopping graphql service")
	}
	log.Info("graphql service stopped")
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
