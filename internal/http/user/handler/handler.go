package handler

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"user_advt/internal/domain/users"
	"user_advt/internal/lib/api/response"
	"user_advt/internal/storage"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

//go:generate mockery --name UserService --output ./mocks --dir .
type UserService interface {
	SignUpUser(ctx context.Context, name, email, password string) (int, error)
	GetUser(ctx context.Context, id string) (users.User, error)
	SignInUser(ctx context.Context, email, password string) (users.User, error)
}

type CreateUserRequestParams struct {
	Name     string `json:"name" binding:"required,min=2,max=20"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type AuthUserRequestParams struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type handler struct {
	logger  *slog.Logger
	service UserService
}

func NewHandler(logger *slog.Logger, service UserService) *handler {
  return &handler{
		logger: logger,
		service: service,
	}
}

func (r *handler) CreateUser(c *gin.Context) {
	var params CreateUserRequestParams

	if err := c.ShouldBindJSON(&params); err != nil {
		c.Error(GetValidError(err))
		r.logger.Error("Failed to bind JSON", slog.String("error:", GetValidError(err).ValidErr))
		return
	}

	r.logger.Info("Binded JSON successfully", slog.String("request:", fmt.Sprintf("%+v", params)))

	userID, err := r.service.SignUpUser(c.Request.Context(), params.Name, params.Email, params.Password)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, map[string]any{
		"success": true,
	  "user_id":  userID,
	})
}	

func (r *handler) GetUserById(c *gin.Context) {
	id := c.Param("id")

	if _, err := strconv.Atoi(id); err != nil {
		c.Error(storage.ErrUserIDInvalidParams)
		return
	}

	user, err := r.service.GetUser(c.Request.Context(), id)

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, map[string]any{
		"success": true,
	  "name": user.Name,
		"email": user.Email,
	})	
}

func (r *handler) AuthUser(c * gin.Context) {
	var params AuthUserRequestParams

	if err := c.ShouldBindJSON(&params); err != nil {
		c.Error(GetValidError(err))
		r.logger.Error("Failed to bind JSON", slog.String("error:", err.Error()))
		return
	}

	user, err := r.service.SignInUser(c.Request.Context(), params.Email, params.Password)
	
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, map[string]any{
		"success": true,
	  "user_id":  user.ID,
	})

}

func GetValidError(err error) *response.ValidError {
	validateError := err.(validator.ValidationErrors)
	return response.ValidationError(validateError)
}

const (
	userUrl = "/users/:id"
	signUpUserUrl = "/users/sign-up"
	singInUserUrl = "/users/sign-in"
)

func (r *handler) Register(router *gin.Engine) {	
  router.POST(signUpUserUrl, r.CreateUser)
  router.POST(singInUserUrl, r.AuthUser)
	router.GET(userUrl, r.GetUserById)
}