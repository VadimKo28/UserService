package main

import (
	"context"
	"fmt"
	// "fmt"
	"log"
	"log/slog"
	"os"
	"user_advt/internal/config"

	// "user_advt/internal/domain/users"
	"user_advt/internal/service"
	"user_advt/internal/storage/user/postgres"
	"user_advt/pkg/hash"
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

	hasher := hash.NewSHA1Hasher("8kWA5D7R)_6@Xv4")

	// fmt.Println(storage.SaveUser(ctx))
	
	service := service.NewUserService(storage, hasher)


	userName, err := service.SignUp(ctx, user.Name, user.Email, user.Password)

	if err != nil {
		logger.Error("Failed to create user", slog.String("error:", err.Error()))
		os.Exit(1)
	}
	fmt.Println(userName, "created")

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