package app

import (
	"github.com/MiKaMoRe/medical-task-tracker/internal/config"
	"github.com/MiKaMoRe/medical-task-tracker/internal/db/txrunner"
	"github.com/MiKaMoRe/medical-task-tracker/internal/http/handler/task"
	taskHandler "github.com/MiKaMoRe/medical-task-tracker/internal/http/handler/task"
	"github.com/MiKaMoRe/medical-task-tracker/internal/logger"
	taskRepository "github.com/MiKaMoRe/medical-task-tracker/internal/repository/task"
	taskService "github.com/MiKaMoRe/medical-task-tracker/internal/service/task"
	"gorm.io/gorm"
)

type App struct {
	cfg    *config.Config
	logger logger.Logger
	db     *gorm.DB

	taskHandler *task.TaskHandler
}

func NewApp(cfg *config.Config, gormDB *gorm.DB, log logger.Logger) *App {
	tx := txrunner.New(gormDB)
	taskRepo := taskRepository.NewTaskRepository(gormDB)
	taskService := taskService.NewTaskService(taskRepo, tx)
	taskHandler := taskHandler.NewTaskHandler(taskService, log)
	return &App{
		cfg:         cfg,
		logger:      log,
		db:          gormDB,
		taskHandler: taskHandler,
	}
}
