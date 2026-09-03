package httputil

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HealthLive(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func HealthReady(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ready"})
}

func RegisterHealthRoutes(r *gin.RouterGroup) {
	r.GET("/health/live", HealthLive)
	r.GET("/health/ready", HealthReady)
}
