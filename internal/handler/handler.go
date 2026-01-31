package handler

import (
	"context"
	"log/slog"
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
	SignInUser(ctx context.Context, email, password string) (string, error)
}

type TokenService interface {
	ParseToken(token string) (string, error)
}

type handler struct {
	logger       *slog.Logger
	service      UserService
	tokenService TokenService
}

func NewHandler(logger *slog.Logger, service UserService, tokenService TokenService) *handler {
	return &handler{
		logger:       logger,
		service:      service,
		tokenService: tokenService,
	}
}

func GetValidError(err error) *response.ValidError {
	validateError := err.(validator.ValidationErrors)
	return response.ValidationError(validateError)
}

const (
	userUrl       = "/users/:id"
	signUpUserUrl = "/users/sign-up"
	singInUserUrl = "/users/sign-in"
)

func (r *handler) Register(router *gin.Engine) {
	router.Use(middleware.ErrorHandler())

	api := router.Group("/api")
	api.Use(Authentication(r))
	api.GET(userUrl, r.GetUserById)

	router.POST(signUpUserUrl, r.SignUp)
	router.POST(singInUserUrl, r.SignIn)
}
