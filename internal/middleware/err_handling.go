package middleware

import (
	"errors"
	"net/http"
	"user_advt/internal/lib/api/response"
	"user_advt/internal/storage"

	"github.com/gin-gonic/gin"
)

// ErrorHandler captures errors and returns a consistent JSON error response
func ErrorHandler() gin.HandlerFunc {
  return func(c *gin.Context) {
    c.Next()
    if len(c.Errors) > 0 {
      err := c.Errors.Last().Err

      if errors.Is(err, storage.ErrUserAlreadyExists) {
        c.JSON(http.StatusBadRequest, map[string]any{
          "success": false,
          "message": err.Error(),
          "status": storage.StatusBadRequest,
        })		
      
        return
      } 

      if errors.Is(err, storage.ErrUserIDInvalidParams) {
        c.JSON(http.StatusBadRequest, map[string]any{
          "success": false,
          "message": err.Error(),
          "status": storage.StatusBadRequest,
        })		
      
        return
      }
      
      if errors.Is(err, storage.ErrUserInvalidCredentials) {
        c.JSON(http.StatusUnauthorized, map[string]any{
          "success": false,
          "message": err.Error(),
          "status": storage.StatusUnauthorized,
        })		
      
        return
      }
        
      if errors.Is(err, storage.ErrUserNotFount) {
        c.JSON(http.StatusNotFound, map[string]any{
          "success": false,
          "message": err.Error(),
          "status": storage.StatusNotFound,
        })		
      
        return
      }

      var validErr *response.ValidError
      
      if errors.As(err, &validErr) {
        c.JSON(http.StatusUnprocessableEntity, map[string]any{
          "success": false,
          "message": validErr.ValidErr,
          "status": validErr.Status,
        })		
      
        return
      }

      c.JSON(http.StatusInternalServerError, map[string]any{
          "success": false,
          "message": err.Error(),
          "status": storage.StatusInternalServerError,
        })	
    }
  }
}