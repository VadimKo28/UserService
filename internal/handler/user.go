package handler

import (
	"app/internal/storage"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *handler) GetUserById(c *gin.Context) {
	id := c.Param("id")

	if _, err := strconv.Atoi(id); err != nil {
		c.Error(storage.ErrUserIDInvalidParams)
		return
	}

	user, err := h.service.GetUser(c.Request.Context(), id)

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, map[string]any{
		"success": true,
		"name":    user.Name,
		"email":   user.Email,
	})
}
