package app

import (
	"net/http"

	"github.com/MiKaMoRe/medical-task-tracker/internal/config"
	"github.com/MiKaMoRe/medical-task-tracker/internal/logger"
	"gorm.io/gorm"
)

type App struct {
	cfg    *config.Config
	logger logger.Logger
	routes routes
}

type routes struct {
	// TODO: Add routes here
}

func NewApp(cfg *config.Config, gormDB *gorm.DB, log logger.Logger) *App {

	return &App{
		cfg:    cfg,
		logger: log,
		routes: routes{},
	}
}

func (a *App) RegisterRoutes(r *http.ServeMux) {
	// TODO: Add routes here
}
