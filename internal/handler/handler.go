package handler

import (
	"context"
	"log/slog"
	"user_advt/internal/domain/subscription"
	"user_advt/internal/domain/users"
	"user_advt/internal/lib/api/response"
	"user_advt/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

//go:generate mockery --name UserService --output ./mocks --dir .
type UserService interface {
	SignUpUser(ctx context.Context, name, email, password string) (int, error)
	GetUser(ctx context.Context, id string) (users.User, error)
	SignInUser(ctx context.Context, email, password string) (string, string, error)
	RefreshTokens(ctx context.Context, refreshToken string) (string, string, error)
	LogOut(ctx context.Context, refreshToken string) error
}

//go:generate mockery --name TokenService --output ./mocks --dir .
type TokenService interface {
	ParseToken(token string) (string, error)
}

//go:generate mockery --name SubscriptionService --output ./mocks --dir .
type SubscriptionService interface {
	CreateSubscription(ctx context.Context, subscriptionDTO subscription.CreateSubscriptionDTO) (int, error)
	GetSubscriptionsByUserID(ctx context.Context, userID, limit, offset int) ([]subscription.Subscription, error)
	UpdateSubscription(ctx context.Context, subscriptionDTO subscription.UpdateSubscriptionDTO) error
}

type handler struct {
	logger              *slog.Logger
	service             UserService
	tokenService        TokenService
	subscriptionService SubscriptionService
}

func NewHandler(logger *slog.Logger, service UserService, tokenService TokenService, subscriptionService SubscriptionService) *handler {
	return &handler{
		logger:              logger,
		service:             service,
		tokenService:        tokenService,
		subscriptionService: subscriptionService,
	}
}

func GetValidError(err error) *response.ValidError {
	validateError := err.(validator.ValidationErrors)
	return response.ValidationError(validateError)
}

const (
	userUrl               = "/users/:id"
	signUpUserUrl         = "/users/sign-up"
	singInUserUrl         = "/users/sign-in"
	refreshUrl            = "/refresh"
	logOutUrl             = "/logout"
	createSubscriptionUrl = "/users/:id/subscriptions"
	updateSubscriptionUrl = "/users/:id/subscriptions/:subscription_id"
)

func (h *handler) Register(router *gin.Engine) {
	router.Use(middleware.ErrorHandler())

	api := router.Group("/api")
	api.Use(Authentication(h))
	api.GET(userUrl, h.GetUserById)
	api.GET(createSubscriptionUrl, h.GetUserSubscriptions)
	api.POST(createSubscriptionUrl, h.CreateUserSubscription)
	api.PUT(updateSubscriptionUrl, h.UpdateUserSubscription)

	router.POST(signUpUserUrl, h.SignUp)
	router.POST(singInUserUrl, h.SignIn)
	router.POST(refreshUrl, h.Refresh)
	router.POST(logOutUrl, h.LogOut)
}
