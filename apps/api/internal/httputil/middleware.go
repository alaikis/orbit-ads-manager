package httputil

import (
	"net/http"

	"orbit/apps/api/config"

	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(error); ok {
			InternalError(c, err.Error())
		} else {
			InternalError(c, "internal server error")
		}
		c.Abort()
	})
}

func CORSMiddleware() gin.HandlerFunc {
	allowedOrigins := config.AppCfg.CORS.AllowedOrigins

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if len(allowedOrigins) > 0 && origin != "" {
			for _, allowed := range allowedOrigins {
				if origin == allowed {
					c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
					break
				}
			}
		}
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, Idempotency-Key")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
