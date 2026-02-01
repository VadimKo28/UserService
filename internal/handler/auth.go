package handler

import (
	"fmt"
	"log/slog"

	"github.com/gin-gonic/gin"
)

type SignInParams struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type SignUpParams struct {
	Name     string `json:"name" binding:"required,min=2,max=20"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func (r *handler) SignUp(c *gin.Context) {
	var params SignUpParams

	if err := c.ShouldBindJSON(&params); err != nil {
		c.Error(GetValidError(err))
		r.logger.Error("Failed to bind JSON", slog.String("error:", GetValidError(err).ValidErr))
		return
	}

	r.logger.Info("Binded JSON successfully", slog.String("request:", fmt.Sprintf("%+v", params)))

	userID, err := r.service.SignUpUser(c.Request.Context(), params.Name, params.Email, params.Password)
	if err != nil {
		c.Error(err)
		r.logger.Error("Failed to sign up user", slog.String("error:", err.Error()))
		return
	}

	c.JSON(200, map[string]any{
		"success": true,
	  "user_id":  userID,
	})
}	


func (r *handler) SignIn(c *gin.Context) {
	var params SignInParams

	if err := c.ShouldBindJSON(&params); err != nil {
		c.Error(GetValidError(err))
		r.logger.Error("Failed to bind JSON", slog.String("error:", err.Error()))
		return
	}

	token, err := r.service.SignInUser(c.Request.Context(), params.Email, params.Password)
	
	if err != nil {
		c.Error(err)
		r.logger.Error("Failed to sign in user", slog.String("error:", err.Error()))
		return
	}

	c.JSON(200, map[string]any{
		"success": true,
	  "token":  token,
	})

}

func (r *handler) Refresh(c *gin.Context) {
	
}
