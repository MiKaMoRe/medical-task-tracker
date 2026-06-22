package middleware

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/MiKaMoRe/medical-task-tracker/internal/httpctx"
)

type InfoLogger interface {
	Info(msg string, args ...any)
}

type LoggerMiddleware struct {
	logger InfoLogger
}

func NewLoggerMiddleware(logger InfoLogger) *LoggerMiddleware {
	return &LoggerMiddleware{logger}
}

type statusResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *statusResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (m *LoggerMiddleware) ReqLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		queryParams := r.URL.Query()
		path := r.URL.Path

		var requestBody []byte
		if r.Body != nil {
			var err error
			requestBody, err = io.ReadAll(r.Body)
			if err == nil {
				r.Body = io.NopCloser(bytes.NewBuffer(requestBody))
			}
		}

		wrappedWriter := &statusResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrappedWriter, r)

		end := time.Now()
		latency := end.Sub(start)

		ctx := r.Context()
		key, ok := httpctx.RequestID(ctx)
		if !ok {
			key = "unknown"
		}

		message := fmt.Sprintf(
			"%s::%s %d %s",
			r.Method,
			path,
			wrappedWriter.statusCode,
			latency,
		)
		m.logger.Info(
			message,
			"queryParams", queryParams,
			"body", string(requestBody),
			"requestId", key,
		)
	})
}
