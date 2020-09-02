package health

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
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

func New(port int, checks []Check) *Server {
	service := &Server{
		http: &http.Server{
			Addr: fmt.Sprintf(":%d", port),
		},
		checks: checks,
	}

	service.setupHandlers()

	return service
}

func (s *Server) setupHandlers() {
	handler := http.NewServeMux()

	handler.HandleFunc("/health", s.serve)

	s.http.Handler = handler
}

func (s *Server) Run(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	log.Info("health check service: starting")

	go func() {
		defer wg.Done()
		log.Debug("health check service: addr=", s.http.Addr)
		err := s.http.ListenAndServe()
		s.runErr = err
		log.Info("health check service: end running > ", err)
	}()

	go func() {
		<-ctx.Done()
		sdCtx, _ := context.WithTimeout(context.Background(), 5*time.Second) // nolint
		err := s.http.Shutdown(sdCtx)
		if err != nil {
			log.Info("health check service shutdown (", err, ")")
		}
	}()

	s.readiness = true
}

func (s *Server) serve(w http.ResponseWriter, r *http.Request) {
	errs := make([]error, 0)
	for _, check := range s.checks {
		if err := check(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		w.WriteHeader(http.StatusInternalServerError)

		for _, err := range errs {
			log.Errorf("health check received error: %s", err.Error())

			_, errWrite := w.Write([]byte(fmt.Sprintf("%s\n", err.Error())))
			if errWrite != nil {
				log.Errorf("health check response write error: %s", errWrite.Error())
			}
		}

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
