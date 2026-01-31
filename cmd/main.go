package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"user_advt/internal/config"
	"user_advt/internal/handler"
	"user_advt/internal/lib/logger"
	"user_advt/internal/service"
	"user_advt/internal/storage/user/postgres"
	"user_advt/pkg/hash"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.MustLoad()
	log.Println("Config loaded")

	logger := logger.LoggerSetup(cfg.Env)
	logger.Info("Init app", slog.String("env:", cfg.Env))

	ctx := context.Background()
	storage, err := postgres.NewStorage(ctx, &cfg, logger)

	if err != nil {
		logger.Error("Failed to init storage", slog.String("error:", err.Error()))
		os.Exit(1)
	}

	hasher := hash.NewSHA1Hasher(cfg.Auth.PasswordSalt)

	tokenService := service.NewTokenService([]byte(cfg.Auth.JWTSecret), cfg.Auth.TokenTTL)
	service := service.NewUserService(storage, hasher, tokenService)

	router := gin.Default()

	handler := handler.NewHandler(logger, service, tokenService)

	handler.Register(router)

	logger.Info("Starting server on port: " + cfg.HTTPServer.Port)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.HTTPServer.Port),
		Handler:      router.Handler(),
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		logger.Error("Failed to start server", slog.String("error:", err.Error()))
		os.Exit(1)
	}

}
