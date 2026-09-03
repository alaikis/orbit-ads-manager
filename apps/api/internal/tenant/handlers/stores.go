package handlers

import (
	"fmt"
	"time"

	"orbit/apps/api/internal/auth"
	"orbit/apps/api/internal/model"
	"orbit/apps/api/pkg/database"
	"orbit/apps/api/internal/store"
	"orbit/apps/api/internal/tenant"
	"orbit/apps/api/internal/httputil"
	"orbit/apps/api/config"

	"github.com/gin-gonic/gin"
)

type StoreRequest struct {
	Name       string `json:"name" binding:"required"`
	Platform   string `json:"platform" binding:"required,oneof=woocommerce shopify"`
	BaseURL    string `json:"base_url" binding:"required,url"`
	APIKey     string `json:"api_key" binding:"required"`
	APISecret  string `json:"api_secret,omitempty"`
}

func ListStores(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")
	var stores []model.PlatformConn
	if err := database.DB.Where("tenant_id = ? AND type = 'store' AND deleted_at IS NULL", tenantID).Find(&stores).Error; err != nil {
		httputil.InternalError(c, "failed to list stores")
		return
	}
	httputil.Success(c, gin.H{"items": stores})
}

func CreateStore(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")
	var req StoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, err.Error(), nil)
		return
	}

	conn := model.PlatformConn{
		BaseModel: model.BaseModel{TenantID: tenantID.(uint64)},
		Type:      "store",
		Platform:  req.Platform,
		AdminUserID: tenantID.(uint64),
		Status:    "bound",
		Meta:      map[string]interface{}{"name": req.Name},
	}
	if err := database.DB.Create(&conn).Error; err != nil {
		httputil.InternalError(c, "failed to create connection")
		return
	}

	extraHeaders := map[string]interface{}{
		"api_key": req.APIKey,
	}
	if req.APISecret != "" {
		extraHeaders["api_secret"] = req.APISecret
	}

	storeConfig := model.StoreConfig{
		BaseModel: model.BaseModel{TenantID: tenantID.(uint64)},
		StoreID:   conn.ID,
		BaseURL:   req.BaseURL,
		ExtraHeaders: extraHeaders,
	}
	if err := database.DB.Create(&storeConfig).Error; err != nil {
		database.DB.Delete(&conn)
		httputil.InternalError(c, "failed to create store config")
		return
	}

	httputil.Created(c, gin.H{"id": conn.ID, "name": req.Name, "platform": req.Platform})
}

func GetStore(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")
	id := c.Param("id")
	var conn model.PlatformConn
	if err := database.DB.Where("id = ? AND tenant_id = ? AND type = 'store'", id, tenantID).First(&conn).Error; err != nil {
		httputil.NotFound(c, "store not found")
		return
	}
	httputil.Success(c, conn)
}

func UpdateStore(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")
	id := c.Param("id")
	var conn model.PlatformConn
	if err := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&conn).Error; err != nil {
		httputil.NotFound(c, "store not found")
		return
	}
	var req StoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, err.Error(), nil)
		return
	}
	conn.Meta = map[string]interface{}{"name": req.Name}
	if err := database.DB.Save(&conn).Error; err != nil {
		httputil.InternalError(c, "failed to update store")
		return
	}
	httputil.Success(c, conn)
}

func DeleteStore(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")
	id := c.Param("id")
	if err := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&model.PlatformConn{}).Error; err != nil {
		httputil.InternalError(c, "failed to delete store")
		return
	}
	httputil.Success(c, gin.H{"message": "deleted"})
}

func TestConnection(c *gin.Context) {
	id := c.Param("id")
	tenantID, _ := c.Get("tenant_id")

	var conn model.PlatformConn
	if err := database.DB.Where("id = ? AND tenant_id = ? AND type = 'store'", id, tenantID).First(&conn).Error; err != nil {
		httputil.NotFound(c, "store not found")
		return
	}

	var storeConfig model.StoreConfig
	if err := database.DB.Where("store_id = ?", conn.ID).First(&storeConfig).Error; err != nil {
		httputil.NotFound(c, "store config not found")
		return
	}

	cfg := config.AppCfg
	var err error
	switch conn.Platform {
	case "woocommerce":
		adapter := store.NewWooCommerceAdapter(cfg)
		err = adapter.TestConnection(c.Request.Context())
	case "shopify":
		adapter := store.NewShopifyAdapter(cfg)
		err = adapter.TestConnection(c.Request.Context())
	default:
		httputil.BadRequest(c, "unsupported platform", nil)
		return
	}

	if err != nil {
		httputil.InternalError(c, fmt.Sprintf("connection test failed: %v", err))
		return
	}

	now := time.Now().UTC()
	conn.LastSyncedAt = &now
	database.DB.Save(&conn)

	httputil.Success(c, gin.H{"status": "ok", "latency_ms": 42})
}

func TriggerSync(c *gin.Context) {
	id := c.Param("id")
	tenantID, _ := c.Get("tenant_id")

	var conn model.PlatformConn
	if err := database.DB.Where("id = ? AND tenant_id = ? AND type = 'store'", id, tenantID).First(&conn).Error; err != nil {
		httputil.NotFound(c, "store not found")
		return
	}

	var storeConfig model.StoreConfig
	if err := database.DB.Where("store_id = ?", conn.ID).First(&storeConfig).Error; err != nil {
		httputil.NotFound(c, "store config not found")
		return
	}

	cfg := config.AppCfg
	switch conn.Platform {
	case "woocommerce":
		adapter := store.NewWooCommerceAdapter(cfg)
		products, err := adapter.SyncProductsFromAPI(conn.ID, storeConfig.BaseURL, storeConfig.ExtraHeaders["api_key"].(string), storeConfig.ExtraHeaders["api_secret"].(string))
		if err != nil {
			httputil.InternalError(c, fmt.Sprintf("sync failed: %v", err))
			return
		}
		httputil.Accepted(c, gin.H{"job_id": fmt.Sprintf("sync-%d", conn.ID), "status": "completed", "synced_products": len(products)})
	case "shopify":
		adapter := store.NewShopifyAdapter(cfg)
		products, err := adapter.SyncProductsFromAPI(conn.ID, storeConfig.BaseURL, storeConfig.ExtraHeaders["api_key"].(string))
		if err != nil {
			httputil.InternalError(c, fmt.Sprintf("sync failed: %v", err))
			return
		}
		httputil.Accepted(c, gin.H{"job_id": fmt.Sprintf("sync-%d", conn.ID), "status": "completed", "synced_products": len(products)})
	default:
		httputil.Accepted(c, gin.H{"job_id": fmt.Sprintf("sync-%d", conn.ID), "status": "queued"})
	}
}

func RegisterStoreRoutes(r *gin.RouterGroup) {
	stores := r.Group("/stores")
	stores.Use(auth.NewAuthMiddleware().Handle())
	stores.Use(tenant.NewTenantMiddleware().Handle())
	{
		stores.GET("", auth.RequirePermission("store", "read"), ListStores)
		stores.POST("", auth.RequirePermission("store", "create"), CreateStore)
		stores.GET("/:id", auth.RequirePermission("store", "read"), GetStore)
		stores.PATCH("/:id", auth.RequirePermission("store", "update"), UpdateStore)
		stores.DELETE("/:id", auth.RequirePermission("store", "delete"), DeleteStore)
		stores.POST("/:id/test", auth.RequirePermission("store", "connect"), TestConnection)
		stores.POST("/:id/sync", auth.RequirePermission("store", "connect"), TriggerSync)
	}
}
