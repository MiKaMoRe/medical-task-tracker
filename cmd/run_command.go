package main

import (
	"context"
	"errors"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/MiKaMoRe/medical-task-tracker/internal/app"
	"github.com/MiKaMoRe/medical-task-tracker/internal/config"
	"github.com/MiKaMoRe/medical-task-tracker/internal/db/gormdb"
	"github.com/MiKaMoRe/medical-task-tracker/internal/logger"

	"gorm.io/gorm"
)

func runCommand() {
	var log logger.Logger
	cfg, err := config.NewConfig()
	if err != nil {
		log = logger.MustNewWithConfigLevel(config.EnvProd, config.LogLevelError)
		log.Error("Failed to load config", "error", err)
		return
	}
	log = logger.MustNewWithConfigLevel(cfg.Env, cfg.LogLevel)
	defer func() {
		if closeErr := log.Close(); closeErr != nil {
			log.Error("Failed to close logger", "error", closeErr)
		}
	}()
	log.Info("Starting application...")

	gormDB, err := setupDB(cfg, log)
	if err != nil {
		log.Error("Error while connect to Database", "error", err)
		return
	}
	log.Info(
		"Database connected successfully",
		"host", cfg.DB.Host,
		"database", cfg.DB.Db,
	)

	application := app.NewApp(cfg, gormDB, log)

	srv := &http.Server{
		Addr:              ":" + cfg.AppPort,
		Handler:           application.Handler(),
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	serverErr := make(chan error, 1)
	go func() {
		log.Info("Server listening", "port", cfg.AppPort)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("Error", "error", err)
			serverErr <- err
		}
	}()

	// Graceful shutdown
	signalCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	select {
	case <-signalCtx.Done():
		log.Info("Shutdown signal received")
	case err := <-serverErr:
		log.Error("Server error, initializing shutdown", "error", err)
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer shutdownCancel()

	log.Info("Shutting down HTTP server...")
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("HTTP server forced to shutdown", "error", err)
		if closeErr := srv.Close(); closeErr != nil {
			log.Error("Failed to close HTTP server", "error", closeErr)
		}
	} else {
		log.Info("HTTP server stopped")
	}

	log.Info("Closing database connection...")
	if err := gormdb.Close(gormDB); err != nil {
		log.Error("Error closing database connection", "error", err)
	}

	log.Info("Server exited gracefully")
}

func setupDB(cfg *config.Config, log logger.Logger) (*gorm.DB, error) {
	dialector := gormdb.NewPostgresDialector(cfg.DB.PostgresDSN())

	opts := &gormdb.Options{
		LogLevel:  mapConfigLogLevelToGorm(cfg.LogLevel),
		AppLogger: log,
	}

	return gormdb.Connect(dialector, opts)
}

func mapConfigLogLevelToGorm(level string) gormdb.LogLevel {
	switch level {
	case config.LogLevelDebug, config.LogLevelInfo:
		return gormdb.LogLevelInfo
	case config.LogLevelWarn:
		return gormdb.LogLevelWarn
	case config.LogLevelError:
		return gormdb.LogLevelError
	default:
		return gormdb.LogLevelError
	}
}
