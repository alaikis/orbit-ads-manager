package handlers

import (
	"orbit/apps/api/internal/auth"
	"orbit/apps/api/internal/model"
	"orbit/apps/api/pkg/database"
	"orbit/apps/api/internal/tenant"
	"orbit/apps/api/internal/httputil"

	"github.com/gin-gonic/gin"
)

func ListRules(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")
	var rules []model.Rule
	if err := database.DB.Where("tenant_id = ? AND deleted_at IS NULL", tenantID).Order("created_at DESC").Find(&rules).Error; err != nil {
		httputil.InternalError(c, "failed to list rules")
		return
	}
	httputil.Success(c, gin.H{"items": rules})
}

func CreateRule(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")
	userID, _ := c.Get("user_id")
	var rule model.Rule
	if err := c.ShouldBindJSON(&rule); err != nil {
		httputil.BadRequest(c, err.Error(), nil)
		return
	}
	rule.TenantID = tenantID.(uint64)
	rule.CreatedBy = userID.(uint64)
	if err := database.DB.Create(&rule).Error; err != nil {
		httputil.InternalError(c, "failed to create rule")
		return
	}
	httputil.Created(c, rule)
}

func GetRule(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")
	id := c.Param("id")
	var rule model.Rule
	if err := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&rule).Error; err != nil {
		httputil.NotFound(c, "rule not found")
		return
	}
	httputil.Success(c, rule)
}

func UpdateRule(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")
	id := c.Param("id")
	var rule model.Rule
	if err := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&rule).Error; err != nil {
		httputil.NotFound(c, "rule not found")
		return
	}
	if err := c.ShouldBindJSON(&rule); err != nil {
		httputil.BadRequest(c, err.Error(), nil)
		return
	}
	if err := database.DB.Save(&rule).Error; err != nil {
		httputil.InternalError(c, "failed to update rule")
		return
	}
	httputil.Success(c, rule)
}

func DeleteRule(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")
	id := c.Param("id")
	if err := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&model.Rule{}).Error; err != nil {
		httputil.InternalError(c, "failed to delete rule")
		return
	}
	httputil.Success(c, gin.H{"message": "deleted"})
}

func ToggleRule(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")
	id := c.Param("id")
	enable := c.Param("action") == "enable"
	var rule model.Rule
	if err := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&rule).Error; err != nil {
		httputil.NotFound(c, "rule not found")
		return
	}
	rule.Enabled = enable
	if err := database.DB.Save(&rule).Error; err != nil {
		httputil.InternalError(c, "failed to toggle rule")
		return
	}
	httputil.Success(c, gin.H{"message": "rule updated", "enabled": enable})
}

func RegisterRuleRoutes(r *gin.RouterGroup) {
	rules := r.Group("/rules")
	rules.Use(auth.NewAuthMiddleware().Handle())
	rules.Use(tenant.NewTenantMiddleware().Handle())
	{
		rules.GET("", auth.RequirePermission("rule", "read"), ListRules)
		rules.POST("", auth.RequirePermission("rule", "create"), CreateRule)
		rules.GET("/:id", auth.RequirePermission("rule", "read"), GetRule)
		rules.PATCH("/:id", auth.RequirePermission("rule", "update"), UpdateRule)
		rules.DELETE("/:id", auth.RequirePermission("rule", "delete"), DeleteRule)
		rules.POST("/:id/enable", auth.RequirePermission("rule", "enable"), func(c *gin.Context) { ToggleRule(c) })
		rules.POST("/:id/disable", auth.RequirePermission("rule", "disable"), func(c *gin.Context) { ToggleRule(c) })
	}
}
