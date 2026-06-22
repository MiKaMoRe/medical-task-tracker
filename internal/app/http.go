package app

import (
	"net/http"

	"github.com/MiKaMoRe/medical-task-tracker/internal/http/middleware"
)

func (a *App) Handler() http.Handler {
	mux := http.NewServeMux()
	a.registerRoutes(mux)

	return chainMiddleware(
		mux,
		middleware.NewLoggerMiddleware(a.logger).ReqLoggerMiddleware,
		middleware.ReqIDMiddleware,
	)
}

func (a *App) registerRoutes(mux *http.ServeMux) {
}

func chainMiddleware(root http.Handler, mws ...func(http.Handler) http.Handler) http.Handler {
	for _, mw := range mws {
		root = mw(root)
	}
	return root
}
