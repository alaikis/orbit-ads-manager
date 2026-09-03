package handlers

import (
	"net/http"

	"orbit/apps/api/internal/auth"
	"orbit/apps/api/internal/model"
	"orbit/apps/api/pkg/database"
	"orbit/apps/api/internal/tenant"
	"orbit/apps/api/internal/httputil"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ListFeeds(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")
	var feeds []model.Feed
	if err := database.DB.Where("tenant_id = ? AND deleted_at IS NULL", tenantID).Find(&feeds).Error; err != nil {
		httputil.InternalError(c, "failed to list feeds")
		return
	}
	httputil.Success(c, gin.H{"items": feeds})
}

func CreateFeed(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")
	var feed model.Feed
	if err := c.ShouldBindJSON(&feed); err != nil {
		httputil.BadRequest(c, err.Error(), nil)
		return
	}
	feed.TenantID = tenantID.(uint64)
	feed.PublicToken = generateToken()
	if err := database.DB.Create(&feed).Error; err != nil {
		httputil.InternalError(c, "failed to create feed")
		return
	}
	httputil.Created(c, feed)
}

func GetFeed(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")
	id := c.Param("id")
	var feed model.Feed
	if err := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&feed).Error; err != nil {
		httputil.NotFound(c, "feed not found")
		return
	}
	httputil.Success(c, feed)
}

func UpdateFeed(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")
	id := c.Param("id")
	var feed model.Feed
	if err := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&feed).Error; err != nil {
		httputil.NotFound(c, "feed not found")
		return
	}
	if err := c.ShouldBindJSON(&feed); err != nil {
		httputil.BadRequest(c, err.Error(), nil)
		return
	}
	if err := database.DB.Save(&feed).Error; err != nil {
		httputil.InternalError(c, "failed to update feed")
		return
	}
	httputil.Success(c, feed)
}

func DeleteFeed(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")
	id := c.Param("id")
	if err := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&model.Feed{}).Error; err != nil {
		httputil.InternalError(c, "failed to delete feed")
		return
	}
	httputil.Success(c, gin.H{"message": "deleted"})
}

func RegenerateFeed(c *gin.Context) {
	httputil.Accepted(c, gin.H{"job_id": "feed-regen-123", "status": "queued"})
}

func GetFeedLogs(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")
	id := c.Param("id")
	var feed model.Feed
	if err := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&feed).Error; err != nil {
		httputil.NotFound(c, "feed not found")
		return
	}
	var runs []model.FeedRun
	if err := database.DB.Where("feed_id = ?", feed.ID).Order("created_at DESC").Limit(20).Find(&runs).Error; err != nil {
		httputil.InternalError(c, "failed to fetch logs")
		return
	}
	httputil.Success(c, gin.H{"items": runs})
}

func GetPublicFeed(c *gin.Context) {
	token := c.Param("token")
	var feed model.Feed
	if err := database.DB.Where("public_token = ?", token).First(&feed).Error; err != nil {
		httputil.NotFound(c, "feed not found")
		return
	}
	c.String(http.StatusOK, "id\ttitle\tprice\tavailability\n")
}

func generateToken() string {
	return "token_" + uuid.New().String()[:16]
}

func RegisterFeedRoutes(r *gin.RouterGroup) {
	feeds := r.Group("/feeds")
	feeds.Use(auth.NewAuthMiddleware().Handle())
	feeds.Use(tenant.NewTenantMiddleware().Handle())
	{
		feeds.GET("", auth.RequirePermission("feed", "read"), ListFeeds)
		feeds.POST("", auth.RequirePermission("feed", "create"), CreateFeed)
		feeds.GET("/:id", auth.RequirePermission("feed", "read"), GetFeed)
		feeds.PATCH("/:id", auth.RequirePermission("feed", "update"), UpdateFeed)
		feeds.DELETE("/:id", auth.RequirePermission("feed", "delete"), DeleteFeed)
		feeds.POST("/:id/regenerate", auth.RequirePermission("feed", "regenerate"), RegenerateFeed)
		feeds.GET("/:id/logs", auth.RequirePermission("feed", "read"), GetFeedLogs)
	}
	publicFeeds := r.Group("/public/feeds")
	{
		publicFeeds.GET("/:token", GetPublicFeed)
	}
}
