package server

import (
	"fmt"
	"github.com/PatrickKvartsaniy/print-me-at/errors"
	"github.com/PatrickKvartsaniy/print-me-at/models"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func (s Server) messageScheduler(w http.ResponseWriter, r *http.Request) {
	msg, err := parseQueryParams(r)
	if err != nil {
		logrus.WithError(err).Error("parsing query params")
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	if err := s.repo.AddNewTask(msg); err != nil {
		logrus.WithError(err).Error("adding new message to the queue")
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func parseQueryParams(r *http.Request) (models.Message, error) {
	msg := r.URL.Query().Get("msg")
	if len(msg) == 0 {
		return models.Message{}, fmt.Errorf("%w: message is missing", errors.InvalidParameters)
	}
	ts := r.URL.Query().Get("ts")
	if len(ts) == 0 {
		return models.Message{}, fmt.Errorf("%w: time is missing", errors.InvalidParameters)
	}
	formattedTime, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		return models.Message{}, err
	}
	return models.Message{
		ID:        uuid.NewV4().String(),
		Value:     msg,
		Timestamp: formattedTime.Unix(),
	}, nil
}
