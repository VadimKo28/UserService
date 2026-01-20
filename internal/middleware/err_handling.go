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

      if err != nil {
        if errors.Is(err, storage.ErrUserAlreadyExists) {
          c.JSON(http.StatusBadRequest, map[string]any{
            "success": false,
            "message": err.Error(),
            "status": storage.StatusBadRequest,
          })		
        
          return
        }  

        var validErr *response.ValidError
        
        if errors.As(err, &validErr) {
          c.JSON(http.StatusUnprocessableEntity, map[string]any{
            "success": false,
            "message": validErr.Err,
            "status": validErr.Status,
          })		
        
          return
        }

      }

      c.JSON(http.StatusInternalServerError, map[string]any{
          "success": false,
          "message": err.Error(),
          "status": storage.StatusInternalServerError,
        })	
    }
  }
}