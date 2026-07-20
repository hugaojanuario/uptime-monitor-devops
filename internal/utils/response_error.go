package utils

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hugaojanuario/uptime-monitor-devops/internal/services"
)

func RespondError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, services.ErrURLNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, services.ErrInvalidURL):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, services.ErrURLAlreadyExists):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
