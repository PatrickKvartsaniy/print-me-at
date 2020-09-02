package health

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

type (
	Server struct {
		http      *http.Server
		checks    []Check
		runErr    error
		readiness bool
	}

	Check func() error
)

func CreateAndRun(port int, checks []Check) *Server {
	service := &Server{
		http: &http.Server{
			Addr: fmt.Sprintf(":%d", port),
		},
		checks: checks,
	}

	service.setupHandlers()
	service.run()

	return service
}

func (s *Server) setupHandlers() {
	handler := http.NewServeMux()

	handler.HandleFunc("/health", s.serve)

	s.http.Handler = handler
}

func (s *Server) run() {
	logrus.Info("health check starting")

	go func() {
		logrus.Debug("health check service: addr=", s.http.Addr)
		if err := s.http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.runErr = err
			logrus.WithError(err).Error("health check service")
		}
	}()

	s.readiness = true
}

func (s *Server) Close(ctx context.Context) {
	if err := s.http.Shutdown(ctx); err != nil {
		logrus.WithError(err).Error("health check service shutdown")
	}
	logrus.Info("health check service stopped")
}

func (s *Server) serve(w http.ResponseWriter, _ *http.Request) {
	errs := make([]error, 0)
	for _, check := range s.checks {
		if err := check(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		w.WriteHeader(http.StatusInternalServerError)

		for _, err := range errs {
			_, errWrite := w.Write([]byte(fmt.Sprintf("%s\n", err.Error())))
			if errWrite != nil {
				logrus.Errorf("health check response write error: %s", errWrite.Error())
			}
		}

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
