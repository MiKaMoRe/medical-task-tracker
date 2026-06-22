package middleware

import (
	"net/http"

	"github.com/MiKaMoRe/medical-task-tracker/internal/httpctx"

	"github.com/google/uuid"
)

func ReqIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := uuid.New().String()
		ctx := httpctx.WithRequestID(r.Context(), reqID)
		w.Header().Set("X-Request-ID", reqID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
