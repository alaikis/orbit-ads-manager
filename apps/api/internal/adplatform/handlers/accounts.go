package handlers

import (
	"orbit/apps/api/internal/auth"
	"orbit/apps/api/internal/model"
	"orbit/apps/api/pkg/database"
	"orbit/apps/api/internal/tenant"
	"orbit/apps/api/internal/httputil"

	"github.com/gin-gonic/gin"
)

func ListAdAccounts(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")
	var accounts []model.AdAccount
	if err := database.DB.Where("tenant_id = ? AND deleted_at IS NULL", tenantID).Find(&accounts).Error; err != nil {
		httputil.InternalError(c, "failed to list ad accounts")
		return
	}
	httputil.Success(c, gin.H{"items": accounts})
}

func CreateAdAccount(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")
	var req struct {
		Name       string `json:"name" binding:"required"`
		Platform   string `json:"platform" binding:"required,oneof=google meta bing tiktok"`
		ExternalID string `json:"external_id" binding:"required"`
		Currency   string `json:"currency"`
		CustomerID string `json:"customer_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, err.Error(), nil)
		return
	}

	account := model.AdAccount{}
	account.TenantID = tenantID.(uint64)
	account.ConnID = 0
	account.Platform = req.Platform
	account.ExternalID = req.ExternalID
	account.Name = req.Name
	account.Currency = req.Currency
	account.CustomerID = req.CustomerID
	account.Status = "active"
	if err := database.DB.Create(&account).Error; err != nil {
		httputil.InternalError(c, "failed to create ad account: "+err.Error())
		return
	}
	httputil.Created(c, account)
}

func GetAdAccount(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")
	id := c.Param("id")
	var account model.AdAccount
	if err := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&account).Error; err != nil {
		httputil.NotFound(c, "ad account not found")
		return
	}
	httputil.Success(c, account)
}

func UpdateAdAccount(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")
	id := c.Param("id")
	var account model.AdAccount
	if err := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&account).Error; err != nil {
		httputil.NotFound(c, "ad account not found")
		return
	}
	httputil.Success(c, account)
}

func DeleteAdAccount(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")
	id := c.Param("id")
	if err := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&model.AdAccount{}).Error; err != nil {
		httputil.InternalError(c, "failed to delete ad account")
		return
	}
	httputil.Success(c, gin.H{"message": "deleted"})
}

func TriggerAdSync(c *gin.Context) {
	httputil.Accepted(c, gin.H{"job_id": "sync-123", "status": "queued"})
}

func RegisterAdAccountRoutes(r *gin.RouterGroup) {
	accounts := r.Group("/ad-accounts")
	accounts.Use(auth.NewAuthMiddleware().Handle())
	accounts.Use(tenant.NewTenantMiddleware().Handle())
	{
		accounts.GET("", auth.RequirePermission("ad_account", "read"), ListAdAccounts)
		accounts.POST("", auth.RequirePermission("ad_account", "connect"), CreateAdAccount)
		accounts.GET("/:id", auth.RequirePermission("ad_account", "read"), GetAdAccount)
		accounts.PATCH("/:id", auth.RequirePermission("ad_account", "update"), UpdateAdAccount)
		accounts.DELETE("/:id", auth.RequirePermission("ad_account", "delete"), DeleteAdAccount)
		accounts.POST("/:id/sync", auth.RequirePermission("ad_account", "update"), TriggerAdSync)
	}
}
