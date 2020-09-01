package server

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/metadata"
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

func correlationID(next http.Handler) http.Handler {
	const key = "X-Correlation-Id"

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		correlationID := uuid.NewV4().String()

		ctx = context.WithValue(ctx, key, correlationID) // nolint
		ctx = metadata.AppendToOutgoingContext(ctx, key, correlationID)

		req = req.WithContext(ctx)
		next.ServeHTTP(w, req)
	})
}
