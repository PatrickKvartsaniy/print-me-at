package server

import (
	"github.com/sirupsen/logrus"
	"net/http"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		logrus.
			WithField("method", req.Method).
			WithField("uri", req.URL.RequestURI()).
			Trace("new http request")

		next.ServeHTTP(w, req)
	})
}
