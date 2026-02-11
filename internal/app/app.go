package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"app/internal/config"
	"app/internal/eventing"
	"app/internal/handler"
	"app/internal/service"
	subscriptionpostgres "app/internal/storage/subscription/postgres"
	tokenpostgres "app/internal/storage/token/postgres"
	userpostgres "app/internal/storage/user/postgres"
	pgclient "app/pkg/client/postgres"
	"app/pkg/hash"

	"github.com/gin-gonic/gin"
)

type Cleanup func() error

func BuildServer(cfg config.Config, logger *slog.Logger) (*http.Server, Cleanup, error) {
	ctx := context.Background()

	pool, err := pgclient.NewClient(ctx, cfg.DatabasePath)
	if err != nil {
		return nil, nil, err
	}

	hasher := hash.NewSHA1Hasher(cfg.Auth.PasswordSalt)

	tokenService := service.NewTokenService([]byte(cfg.Auth.JWTSecret), cfg.Auth.TokenTTL)
	userStorage := userpostgres.NewUserStorage(pool, logger)
	tokenStorage := tokenpostgres.NewTokenStorage(pool, logger)
	subscriptionStorage := subscriptionpostgres.NewSubscriptionStorage(pool, logger)
	userService := service.NewUserService(userStorage, tokenStorage, hasher, tokenService, cfg.Auth.RefreshTokenTTL)
	subscriptionPublisher := eventing.BuildSubscriptionPublisher(cfg, logger)
	subscriptionService := service.NewSubscriptionService(subscriptionStorage, subscriptionPublisher, logger)

	router := gin.Default()
	h := handler.NewHandler(logger, userService, tokenService, subscriptionService)
	h.Register(router)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.HTTPServer.Port),
		Handler:      router.Handler(),
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	cleanup := func() error {
		if subscriptionPublisher == nil {
			return nil
		}
		return subscriptionPublisher.Close()
	}

	return srv, cleanup, nil
}
