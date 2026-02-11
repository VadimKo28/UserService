package main

import (
	"app/internal/app"
	"app/internal/config"
	"app/internal/lib/logger"
	"log"
	"log/slog"
	"os"
)

func main() {
	cfg := config.MustLoad()
	log.Println("Config loaded")

	logger := logger.LoggerSetup(cfg.Env)
	logger.Info("Init app", slog.String("env:", cfg.Env))

	srv, cleanup, err := app.BuildServer(cfg, logger)
	if err != nil {
		logger.Error("Failed to init storage", slog.String("error:", err.Error()))
		os.Exit(1)
	}
	defer func() {
		if cleanup == nil {
			return
		}
		if err := cleanup(); err != nil {
			logger.Error("Failed to close subscription publisher", slog.String("error:", err.Error()))
		}
	}()

	logger.Info("Starting server on port: " + cfg.HTTPServer.Port)

	if err := srv.ListenAndServe(); err != nil {
		logger.Error("Failed to start server", slog.String("error:", err.Error()))
		os.Exit(1)
	}

}
