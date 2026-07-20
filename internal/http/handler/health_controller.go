package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hugaojanuario/uptime-monitor-devops/internal/models"
)

// Health godoc
//
//	@Summary		Healthcheck
//	@Description	Retorna ok quando a api está no ar
//	@Tags			health
//	@Produce		json
//	@Success		200	{object}	models.HealthResponse
//	@Router			/health [get]
func Health(c *gin.Context) {
	c.JSON(http.StatusOK, models.HealthResponse{Status: "ok"})
}
