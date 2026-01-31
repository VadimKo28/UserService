package handler

import (
	"errors"
	"log/slog"
	"strings"
	"user_advt/internal/storage"

	"github.com/gin-gonic/gin"
)

func Authentication(h *handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := getToken(c)

		if err != nil {
			h.logger.Error("Authorization failed", slog.String("error:", err.Error()))
			c.Error(storage.ErrUnauthorized)
			c.Abort()
			return
		}

		userId, err := h.service.ParseToken(token)

		if err != nil {
			h.logger.Error("Authorization failed", slog.String("error:", err.Error()))
			c.Error(storage.ErrUnauthorized)
			c.Abort()
			return
		}

		user, err := h.service.GetUser(c.Request.Context(), userId)
		if err != nil {
			h.logger.Error("Authorization failed", slog.String("error:", err.Error()))
			c.Error(err)
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Set("userId", userId)

		c.Next()
	}
}

func getToken(c *gin.Context) (string, error) {
	header := c.Request.Header.Get("Authorization")

	if header == "" {
		return "", errors.New("Authorization header is empty")
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", errors.New("header is invalid")
	}

	if len(headerParts[1]) == 0 {
		return "", errors.New("token is empty")
	}
	return headerParts[1], nil
}
