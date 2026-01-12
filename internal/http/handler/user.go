package handler

import (
	"fmt"
	"log/slog"
	"os"
	"user_advt/internal/http/lib/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserService interface {
	SignUp(ctx *gin.Context, name, email, password string) (string, error)
	GetUser(ctx *gin.Context, id int) (string, error)
}

type Request struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type Response struct {
	Status string `json:"status"`
	Error string `json:"error,omitempty"`
	Name string `json:"name,omitempty"`	
}

type Handler struct {
	logger  *slog.Logger
	service UserService
}

func NewHandler(logger *slog.Logger, service UserService) *Handler {
  return &Handler{
		logger: logger,
		service: service,
	}
}

func (r *Handler) CreateUser(c *gin.Context) {
	var req Request

	if err := c.ShouldBindJSON(&req); err != nil {
		validateError := err.(validator.ValidationErrors)

		c.JSON(400, Response{
			Status: "error",
			Error:  response.ValidationError(validateError),
		})

		r.logger.Error("Failed to bind JSON", slog.String("error:", response.ValidationError(validateError)))
		return
	}

	r.logger.Info("Binded JSON successfully", slog.String("request:", fmt.Sprintf("%+v", req)))

	userName, err := r.service.SignUp(c, req.Name, req.Email, req.Password)
	if err != nil {
		c.JSON(500, Response{
			Status: "error",
			Error:  err.Error(),
		})
		return
	}

	c.JSON(200, Response{
		Status: "success",
		Name:   userName,
	})
}	

func (r *Handler) InitUserRoutes(router *gin.Engine, port string) {	
  router.POST("/users/sign-up", r.CreateUser)

	r.logger.Info("Starting server on port: " + port )
	if err := router.Run(fmt.Sprintf(":%s", port)); err != nil {
		r.logger.Error("Failed to start server", slog.String("error:", err.Error()))
		os.Exit(1)
	}
}