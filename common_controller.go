package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

const (
	keyLoggerMiddleware = iota
)

func respondWithError(w http.ResponseWriter, req *http.Request, code int, message string) {
	respondWithJSON(w, req, code, map[string]string{"error": message})
}

func headerSetter(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		correlationId := req.Header.Get("X-Correlation-Id")

		if correlationId == "" {
			// For Correlation Id
			correlationId = uuid.Must(uuid.NewV4()).String()
			req.Header.Set("X-Correlation-Id", correlationId)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Correlation-Id", correlationId)

		next(w, req)
	}
}

func (a *App) requestLogger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		reqLogger := log.WithFields(log.Fields{
			"method":    req.Method,
			"timestamp": time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
			"level":     "info",
			"path":      req.URL.Path,
			"message":   req.Method + " " + req.URL.Path})

		reqLogger.Info("Request Received")
		ctx := req.Context()
		ctx = context.WithValue(ctx, keyLoggerMiddleware, reqLogger)
		next(w, req.WithContext(ctx))
	}
}

func respondWithJSON(w http.ResponseWriter, req *http.Request, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	correlationId := w.Header().Get("X-Correlation-Id")

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Correlation-Id", correlationId)
	w.WriteHeader(code)
	w.Write(response)

}
