package main

import (
	"log"
	"log/slog"
	"os"
	"user_advt/internal/config"
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