package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"user_advt/internal/config"
	"user_advt/internal/lib/api/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

//go:generate mockery --name UserService --output ./mocks --dir .
type UserService interface {
	SignUp(ctx *gin.Context, name, email, password string) (int, error)
	GetUser(ctx *gin.Context, id int) (string, error)
}

type Request struct {
	Name     string `json:"name" binding:"required,min=2,max=20"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type Response struct {
	Status string `json:"status"`
	Error string `json:"error,omitempty"`
	UserID int `json:"user_id,omitempty"`	
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

	userID, err := r.service.SignUp(c, req.Name, req.Email, req.Password)
	if err != nil {
		c.JSON(500, Response{
			Status: "error",
			Error:  err.Error(),
		})
		return
	}

	c.JSON(200, Response{
		Status: "success",
	  UserID:  userID,
	})
}	

func (r *Handler) InitUserRouter(router *gin.Engine, cfg config.Config) {	
  router.POST("/users/sign-up", r.CreateUser)

	r.logger.Info("Starting server on port: " + cfg.HTTPServer.Port)

	 srv := &http.Server{
    Addr:    fmt.Sprintf(":%s", cfg.HTTPServer.Port),
    Handler: router.Handler(),
		ReadTimeout: cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout: cfg.HTTPServer.IdleTimeout,
  }

	if err := srv.ListenAndServe(); err != nil {
		r.logger.Error("Failed to start server", slog.String("error:", err.Error()))
		os.Exit(1)
	}
}