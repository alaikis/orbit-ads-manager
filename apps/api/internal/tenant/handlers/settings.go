package handlers

import (
	"orbit/apps/api/internal/auth"
	"orbit/apps/api/internal/model"
	"orbit/apps/api/pkg/database"
	"orbit/apps/api/internal/tenant"
	"orbit/apps/api/internal/httputil"

	"github.com/gin-gonic/gin"
)

func GetTenantSettings(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")
	var tenant model.Tenant
	if err := database.DB.First(&tenant, tenantID).Error; err != nil {
		httputil.NotFound(c, "tenant not found")
		return
	}
	httputil.Success(c, tenant)
}

func UpdateTenantSettings(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")
	var tenant model.Tenant
	if err := database.DB.First(&tenant, tenantID).Error; err != nil {
		httputil.NotFound(c, "tenant not found")
		return
	}
	var req struct {
		Name               string                 `json:"name"`
		Plan               string                 `json:"plan"`
		Timezone           string                 `json:"timezone"`
		AutoApproveHighRisk bool                   `json:"auto_approve_high_risk"`
		Settings           map[string]interface{} `json:"settings"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, err.Error(), nil)
		return
	}
	if req.Name != "" {
		tenant.Name = req.Name
	}
	if req.Plan != "" {
		tenant.Plan = req.Plan
	}
	if req.Timezone != "" {
		tenant.Timezone = req.Timezone
	}
	tenant.AutoApproveHighRisk = req.AutoApproveHighRisk
	if req.Settings != nil {
		tenant.Settings = req.Settings
	}
	if err := database.DB.Save(&tenant).Error; err != nil {
		httputil.InternalError(c, "failed to update tenant")
		return
	}
	httputil.Success(c, tenant)
}

func RegisterTenantSettingsRoutes(r *gin.RouterGroup) {
	tenants := r.Group("/tenants")
	tenants.Use(auth.NewAuthMiddleware().Handle())
	tenants.Use(tenant.NewTenantMiddleware().Handle())
	{
		tenants.GET("/current", auth.RequirePermission("tenant", "read"), GetTenantSettings)
		tenants.PATCH("/current", auth.RequirePermission("tenant", "update"), UpdateTenantSettings)
	}
}
