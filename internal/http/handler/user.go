package handler

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
)

type UserService interface {
	SignUp(ctx *gin.Context, name, email, password string) (string, error)
	GetUser(ctx *gin.Context, id int) (string, error)
}

type Request struct {
	Name     string `json:"name" validate:"required,name"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
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
		c.JSON(400, Response{
			Status: "error",
			Error:  "invalid request body",
		})
		return
	}

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
  router.POST("/users/signup", r.CreateUser)

	r.logger.Info("Starting server on port: " + port )
	if err := router.Run(fmt.Sprintf(":%s", port)); err != nil {
		r.logger.Error("Failed to start server", slog.String("error:", err.Error()))
		os.Exit(1)
	}
}