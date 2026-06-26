package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/MiKaMoRe/medical-task-tracker/internal/httpctx"
)

const (
	maxLoggedRequestBodyBytes = 4096
	redactedPlaceholder       = "***"
)

var sensitiveFields = map[string]struct{}{
	"password":      {},
	"pass":          {},
	"token":         {},
	"access_token":  {},
	"refresh_token": {},
	"secret":        {},
	"authorization": {},
	"api_key":       {},
	"apikey":        {},
}

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

		var requestBodyLogValue string
		if r.Body != nil {
			peekedBody, truncated, err := peekBody(r.Body, maxLoggedRequestBodyBytes)
			if err == nil {
				r.Body = io.NopCloser(io.MultiReader(bytes.NewReader(peekedBody), r.Body))
				requestBodyLogValue = sanitizeBodyForLogs(r.Header.Get("Content-Type"), peekedBody, truncated)
			} else {
				requestBodyLogValue = "<body-read-error>"
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
			"body", requestBodyLogValue,
			"requestId", key,
		)
	})
}

func peekBody(body io.Reader, limit int64) ([]byte, bool, error) {
	reader := io.LimitReader(body, limit+1)
	raw, err := io.ReadAll(reader)
	if err != nil {
		return nil, false, err
	}

	truncated := int64(len(raw)) > limit
	if truncated {
		raw = raw[:limit]
	}
	return raw, truncated, nil
}

func sanitizeBodyForLogs(contentType string, body []byte, truncated bool) string {
	if len(body) == 0 {
		return ""
	}

	ct := strings.ToLower(contentType)
	switch {
	case strings.Contains(ct, "application/json"):
		out := sanitizeJSONBody(body)
		if truncated {
			return out + " [truncated]"
		}
		return out
	case strings.Contains(ct, "application/x-www-form-urlencoded"):
		out := sanitizeFormBody(body)
		if truncated {
			return out + " [truncated]"
		}
		return out
	case strings.HasPrefix(ct, "text/"), strings.Contains(ct, "application/xml"), strings.Contains(ct, "application/yaml"), strings.Contains(ct, "application/javascript"):
		if truncated {
			return string(body) + " [truncated]"
		}
		return string(body)
	default:
		if truncated {
			return "<omitted non-text body> [truncated]"
		}
		return "<omitted non-text body>"
	}
}

func sanitizeJSONBody(body []byte) string {
	var payload any
	if err := json.Unmarshal(body, &payload); err != nil {
		return "<invalid-json>"
	}

	redactSensitiveFields(payload)
	safe, err := json.Marshal(payload)
	if err != nil {
		return "<json-marshal-error>"
	}
	return string(safe)
}

func sanitizeFormBody(body []byte) string {
	values, err := url.ParseQuery(string(body))
	if err != nil {
		return "<invalid-form-body>"
	}

	for key := range values {
		if isSensitiveField(key) {
			values[key] = []string{redactedPlaceholder}
		}
	}
	return values.Encode()
}

func redactSensitiveFields(value any) {
	switch payload := value.(type) {
	case map[string]any:
		for k, v := range payload {
			if isSensitiveField(k) {
				payload[k] = redactedPlaceholder
				continue
			}
			redactSensitiveFields(v)
		}
	case []any:
		for _, item := range payload {
			redactSensitiveFields(item)
		}
	}
}

func isSensitiveField(field string) bool {
	_, ok := sensitiveFields[strings.ToLower(field)]
	return ok
}
