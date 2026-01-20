package handler

import (
	"context"
	"errors"
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
	SignUp(ctx context.Context, name, email, password string) (int, error)
	GetUser(ctx context.Context, id string) (users.GetUserDTO, error)
}

type CreateUserRequestParams struct {
	Name     string `json:"name" binding:"required,min=2,max=20"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type Response struct {
	Status string `json:"status"`
	Error string  `json:"error,omitempty"`
	UserID int `json:"user_id,omitempty"`	
	UserName string `json:"user_name,omitempty"`
	UserEmail string `json:"user_email,omitempty"`
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
		validateError := err.(validator.ValidationErrors)

		c.Error(response.ValidationError(validateError))
		r.logger.Error("Failed to bind JSON", slog.String("error:", response.ValidationError(validateError).Err))
		return
	}

	r.logger.Info("Binded JSON successfully", slog.String("request:", fmt.Sprintf("%+v", params)))

	userID, err := r.service.SignUp(c.Request.Context(), params.Name, params.Email, params.Password)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, Response{
		Status: "success",
	  UserID:  userID,
	})
}	

func (r *handler) GetUserById(c *gin.Context) {

	id := c.Param("id")

	if _, err := strconv.Atoi(id); err != nil {
		c.JSON(400, Response{
			Status: "error",
			Error:  "id params must be an integer",
		})
		return
	}

	user, err := r.service.GetUser(c.Request.Context(), id)

	if err != nil {
		if errors.Is(err, storage.ErrUserNotFount) {
			c.JSON(404, Response{
				Status: "error",
				Error:  storage.ErrUserNotFount.Error(),
			})
			return
	  }
		
		c.JSON(500, Response{
			Status: "error",
			Error:  err.Error(),
		})
		return
	}

	c.JSON(200, Response{
		Status: "success",
		UserID: user.ID,
	  UserName: user.Name,
		UserEmail: user.Email,
	})	
}

const (
	userUrl = "/users/:id"
	signUpUserUrl = "/users/sign-up"
)

func (r *handler) Register(router *gin.Engine) {	
  router.POST(signUpUserUrl, r.CreateUser)
	router.GET(userUrl, r.GetUserById)
}