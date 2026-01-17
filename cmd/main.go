package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"user_advt/internal/config"
	"user_advt/internal/service"
	"user_advt/internal/storage/user/postgres"
	"user_advt/internal/user"
	"user_advt/pkg/hash"

	"github.com/gin-gonic/gin"
)

const (
	envLocal	= "local"
	envProd	= "prod"
	envDev	= "dev"
)

func main() {
	cfg := config.MustLoad()
	log.Println("Config loaded")

	logger := LoggerSetup(cfg.Env)
	logger.Info("Init app", slog.String("env:", cfg.Env))

	ctx := context.Background()
	storage, err := postgres.NewStorage(ctx, &cfg, logger)
	
	if err != nil {
		logger.Error("Failed to init storage", slog.String("error:", err.Error()))
		os.Exit(1)
	}

	hasher := hash.NewSHA1Hasher(os.Getenv("SALT_HASH"))
	
	service := service.NewUserService(storage, hasher)

	router := gin.Default()

	handler := user.NewHandler(logger, service)

	handler.Register(router)

	logger.Info("Starting server on port: " + cfg.HTTPServer.Port)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.HTTPServer.Port),
		Handler: router.Handler(),
		ReadTimeout: cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout: cfg.HTTPServer.IdleTimeout,
  }

	if err := srv.ListenAndServe(); err != nil {
		logger.Error("Failed to start server", slog.String("error:", err.Error()))
		os.Exit(1)
	}

}

func LoggerSetup(env string) *slog.Logger {
	var log *slog.Logger

  switch env {
    case envLocal:
			log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelDebug,
			}))
    case envDev:
			log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelDebug,
			}))
		case envProd:
			log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelInfo,
			}))
  }

	return log 
}