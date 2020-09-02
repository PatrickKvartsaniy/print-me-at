package server

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func (s Server) messageScheduler(w http.ResponseWriter, r *http.Request) {
	msg, ts, err := parseQueryParams(r)
	if err != nil {
		logrus.WithError(err).Error("parsing query params")
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	if err := s.repo.AddNewTask(msg, ts); err != nil {
		logrus.WithError(err).Error("adding new message to the queue")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func parseQueryParams(r *http.Request) (string, time.Time, error) {
	msg := r.URL.Query().Get("msg")
	if len(msg) == 0 {
		return "", time.Time{}, fmt.Errorf("message is missing")
	}
	ts := r.URL.Query().Get("ts")
	if len(ts) == 0 {
		return "", time.Time{}, fmt.Errorf("time is missing")
	}
	formattedTime, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		return "", time.Time{}, err
	}
	return msg, formattedTime, nil
}
