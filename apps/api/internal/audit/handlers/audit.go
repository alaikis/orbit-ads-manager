package handlers

import (
	"orbit/apps/api/internal/auth"
	"orbit/apps/api/internal/model"
	"orbit/apps/api/pkg/database"
	"orbit/apps/api/internal/tenant"
	"orbit/apps/api/internal/httputil"

	"github.com/gin-gonic/gin"
)

func ListAuditLogs(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")
	var logs []model.AuditLog
	if err := database.DB.Where("tenant_id = ?", tenantID).Order("created_at DESC").Limit(100).Find(&logs).Error; err != nil {
		httputil.InternalError(c, "failed to list audit logs")
		return
	}
	httputil.Success(c, gin.H{"items": logs})
}

func RegisterAuditRoutes(r *gin.RouterGroup) {
	audit := r.Group("/audit-logs")
	audit.Use(auth.NewAuthMiddleware().Handle())
	audit.Use(tenant.NewTenantMiddleware().Handle())
	{
		audit.GET("", auth.RequirePermission("audit_log", "read"), ListAuditLogs)
	}
}
