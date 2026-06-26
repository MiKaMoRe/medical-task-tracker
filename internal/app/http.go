package app

import (
	"net/http"

	"github.com/MiKaMoRe/medical-task-tracker/internal/http/middleware"
	"github.com/MiKaMoRe/medical-task-tracker/internal/http/response"
)

func (a *App) Handler() http.Handler {
	mux := http.NewServeMux()
	a.RegisterRoutes(mux)

	return chainMiddleware(
		mux,
		middleware.NewLoggerMiddleware(a.logger).ReqLoggerMiddleware,
		middleware.ReqIDMiddleware,
	)
}

func (a *App) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/tasks/create", a.taskHandler.CreateTask)
	mux.HandleFunc("/api/v1/tasks", a.taskHandler.GetTasks)
	mux.HandleFunc("/api/v1/tasks/{id}", a.taskHandler.TaskByID)
	mux.HandleFunc("/api/v1/tasks/{id}/done", a.taskHandler.MarkTaskDone)
	mux.HandleFunc("/api/v1/tasks/{id}/tags", a.taskHandler.TaskTags)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		response.NotFound(w, "Not Found")
	})
}

func chainMiddleware(root http.Handler, mws ...func(http.Handler) http.Handler) http.Handler {
	for _, mw := range mws {
		root = mw(root)
	}
	return root
}
