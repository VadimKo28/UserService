package handler

import (
	"app/internal/domain/subscription"
	"app/internal/metrics"
	"app/internal/storage"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateSubscriptionParams struct {
	ServiceName string `json:"service_name" binding:"required"`
	Price       int    `json:"price" binding:"required,gt=0"`
	StartDate   string `json:"start_date" binding:"required,datetime=2006-01"`
	EndDate     string `json:"end_date" binding:"omitempty,datetime=2006-01"`
}

type UpdateSubscriptionParams struct {
	ServiceName string `json:"service_name" binding:"required"`
	Price       int    `json:"price" binding:"required,gt=0"`
	StartDate   string `json:"start_date" binding:"required,datetime=2006-01"`
	EndDate     string `json:"end_date" binding:"omitempty,datetime=2006-01"`
}

func (h *handler) CreateUserSubscription(c *gin.Context) {
	var params CreateSubscriptionParams

	if err := c.ShouldBindJSON(&params); err != nil {
		c.Error(GetValidError(err))
		h.logger.Error("Failed to bind JSON", slog.String("error:", GetValidError(err).ValidErr))
		return
	}

	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(storage.ErrUserIDInvalidParams)
		h.logger.Error("Failed to get user id", slog.String("error:", storage.ErrUserIDInvalidParams.Error()))
		return
	}

	startDate, err := time.Parse("2006-01", params.StartDate)
	if err != nil {
		c.Error(fmt.Errorf("invalid start_date format"))
		return
	}

	var endDate time.Time
	if params.EndDate != "" {
		endDate, err = time.Parse("2006-01", params.EndDate)
		if err != nil {
			c.Error(fmt.Errorf("invalid end_date format"))
			return
		}
	}

	subscriptionDTO := subscription.CreateSubscriptionDTO{
		UserID:      userID,
		ServiceName: params.ServiceName,
		Price:       params.Price,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	subscriptionID, err := h.subscriptionService.CreateSubscription(
		c.Request.Context(),
		subscriptionDTO,
	)
	if err != nil {
		c.Error(err)
		h.logger.Error("Failed to create subscription", slog.String("error:", err.Error()))
		return
	}

	metrics.IncSubscriptionsCreated(time.Now())

	c.JSON(http.StatusCreated, map[string]any{
		"success":         true,
		"subscription_id": subscriptionID,
	})
}

func (h *handler) GetUserSubscriptions(c *gin.Context) {
	const (
		defaultLimit = 20
		maxLimit     = 100
	)

	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(storage.ErrUserIDInvalidParams)
		h.logger.Error("Failed to get user id", slog.String("error:", storage.ErrUserIDInvalidParams.Error()))
		return
	}

	limit := defaultLimit
	offset := 0

	if limitParam := c.Query("limit"); limitParam != "" {
		parsedLimit, err := strconv.Atoi(limitParam)
		if err != nil || parsedLimit <= 0 || parsedLimit > maxLimit {
			c.Error(storage.ErrInvalidPaginationParams)
			return
		}
		limit = parsedLimit
	}

	if offsetParam := c.Query("offset"); offsetParam != "" {
		parsedOffset, err := strconv.Atoi(offsetParam)
		if err != nil || parsedOffset < 0 {
			c.Error(storage.ErrInvalidPaginationParams)
			return
		}
		offset = parsedOffset
	}

	subscriptions, err := h.subscriptionService.GetSubscriptionsByUserID(
		c.Request.Context(),
		userID,
		limit,
		offset,
	)
	if err != nil {
		c.Error(err)
		h.logger.Error("Failed to get subscriptions", slog.String("error:", err.Error()))
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"success":       true,
		"limit":         limit,
		"offset":        offset,
		"subscriptions": subscriptions,
	})
}

func (h *handler) UpdateUserSubscription(c *gin.Context) {
	var params UpdateSubscriptionParams

	if err := c.ShouldBindJSON(&params); err != nil {
		c.Error(GetValidError(err))
		h.logger.Error("Failed to bind JSON", slog.String("error:", GetValidError(err).ValidErr))
		return
	}

	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(storage.ErrUserIDInvalidParams)
		h.logger.Error("Failed to get user id", slog.String("error:", storage.ErrUserIDInvalidParams.Error()))
		return
	}

	subscriptionID, err := strconv.Atoi(c.Param("subscription_id"))
	if err != nil {
		c.Error(storage.ErrSubscriptionIDInvalidParams)
		h.logger.Error("Failed to get subscription id", slog.String("error:", storage.ErrSubscriptionIDInvalidParams.Error()))
		return
	}

	startDate, err := time.Parse("2006-01", params.StartDate)
	if err != nil {
		c.Error(fmt.Errorf("invalid start_date format"))
		return
	}

	var endDate time.Time
	if params.EndDate != "" {
		endDate, err = time.Parse("2006-01", params.EndDate)
		if err != nil {
			c.Error(fmt.Errorf("invalid end_date format"))
			return
		}
	}

	subscriptionDTO := subscription.UpdateSubscriptionDTO{
		ID:          subscriptionID,
		UserID:      userID,
		ServiceName: params.ServiceName,
		Price:       params.Price,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	if err := h.subscriptionService.UpdateSubscription(c.Request.Context(), subscriptionDTO); err != nil {
		c.Error(err)
		h.logger.Error("Failed to update subscription", slog.String("error:", err.Error()))
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"success": true,
	})
}
