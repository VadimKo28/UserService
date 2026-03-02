package handler

import (
	"app/internal/domain/subscription"
	"app/internal/domain/users"
	"app/internal/lib/api/response"
	"app/internal/metrics"
	"app/internal/middleware"
	"context"
	"log/slog"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

//go:generate mockery --name UserService --output ./mocks --dir .
type UserService interface {
	SignUpUser(ctx context.Context, name, email, password string) (int, error)
	SignUpUserWithTokens(ctx context.Context, name, email, password string) (int, string, string, error)
	GetUser(ctx context.Context, id string) (users.User, error)
	SignInUser(ctx context.Context, email, password string) (string, string, int, error)
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
	userSubscriptionsUrl  = "/users/:id/subscriptions"
	updateSubscriptionUrl = "/users/:id/subscriptions/:subscription_id"
)

func (h *handler) Register(router *gin.Engine) {
	router.Use(middleware.ErrorHandler())
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://194.156.66.86", "http://localhost:3010", "http://127.0.0.1:3010"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Authorization"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "http://194.156.66.86" ||
				origin == "http://localhost:3010" ||
				origin == "http://127.0.0.1:3010"
		},
		MaxAge: 12 * time.Hour,
	}))

	router.GET("/metrics", gin.WrapH(promhttp.HandlerFor(metrics.Registry, promhttp.HandlerOpts{})))

	api := router.Group("/api")
	public := api.Group("")
	public.POST(signUpUserUrl, h.SignUp)
	public.POST(singInUserUrl, h.SignIn)
	public.POST(refreshUrl, h.Refresh)
	public.POST(logOutUrl, h.LogOut)

	protected := api.Group("")
	protected.Use(Authentication(h))
	protected.GET(userUrl, h.GetUserById)
	protected.GET(userSubscriptionsUrl, h.GetUserSubscriptions)
	protected.POST(userSubscriptionsUrl, h.CreateUserSubscription)
	protected.PUT(updateSubscriptionUrl, h.UpdateUserSubscription)
}
