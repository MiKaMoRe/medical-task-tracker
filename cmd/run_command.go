package main

import (
	"context"
	"errors"
	"net/http"
	"os"
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var log logger.Logger
	cfg, err := config.NewConfig()
	if err != nil {
		log = logger.MustNew(config.EnvProd)
		log.Error("Failed to load config", "error", err)
		return
	}
	log = logger.MustNew(cfg.Env)
	log.Info("Starting application...")

	gormDB, err := setupDB(cfg)
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
		Addr:    ":" + cfg.AppPort,
		Handler: application.Handler(),
	}

	serverErr := make(chan error, 1)
	go func() {
		log.Info("Server listening", "port", cfg.AppPort)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("Error", "error", err)
			serverErr <- err
		}
	}()

	// Gracefull shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(quit)

	select {
	case sig := <-quit:
		log.Info("Shutdown signal received", "signal", sig)
	case err := <-serverErr:
		log.Error("Server error, initializing shutdown", "error", err)
	case <-ctx.Done():
		log.Info("Context cancelled")
	}

	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()

	log.Info("Shutting down HTTP server...")
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("HTTP server forced to shutdown", "error", err)
	} else {
		log.Info("HTTP server stopped")
	}

	log.Info("Closing database connection...")
	if err := gormdb.Close(gormDB); err != nil {
		log.Error("Error closing database connection", "error", err)
	}

	log.Info("Server exited gracefully")
}

func setupDB(cfg *config.Config) (*gorm.DB, error) {
	dialector := gormdb.NewPostgresDialector(cfg.DB.PostgresDSN())

	gormLogLevel := gormdb.LogLevelError
	if cfg.LogLevel == config.LogLevelInfo {
		gormLogLevel = gormdb.LogLevelInfo
	}

	opts := &gormdb.Options{
		LogLevel: gormLogLevel,
	}

	return gormdb.Connect(dialector, opts)
}
