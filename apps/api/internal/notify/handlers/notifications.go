package handlers

import (
	"orbit/apps/api/internal/auth"
	"orbit/apps/api/internal/model"
	"orbit/apps/api/pkg/database"
	"orbit/apps/api/internal/httputil"

	"github.com/gin-gonic/gin"
)

func ListNotifications(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var notifications []model.Notification
	if err := database.DB.Where("user_id = ? AND deleted_at IS NULL", userID).Order("created_at DESC").Limit(50).Find(&notifications).Error; err != nil {
		httputil.InternalError(c, "failed to list notifications")
		return
	}
	httputil.Success(c, gin.H{"items": notifications})
}

func GetUnreadCount(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var count int64
	database.DB.Model(&model.Notification{}).Where("user_id = ? AND read_at IS NULL", userID).Count(&count)
	httputil.Success(c, gin.H{"count": count})
}

func MarkRead(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id := c.Param("id")
	if err := database.DB.Model(&model.Notification{}).Where("id = ? AND user_id = ?", id, userID).Update("read_at", "now()").Error; err != nil {
		httputil.InternalError(c, "failed to mark as read")
		return
	}
	httputil.Success(c, gin.H{"message": "marked as read"})
}

func RegisterNotificationRoutes(r *gin.RouterGroup) {
	notifications := r.Group("/notifications")
	notifications.Use(auth.NewAuthMiddleware().Handle())
	{
		notifications.GET("", auth.RequirePermission("notification", "read"), ListNotifications)
		notifications.GET("/unread-count", auth.RequirePermission("notification", "read"), GetUnreadCount)
		notifications.PATCH("/:id/read", auth.RequirePermission("notification", "mark_read"), MarkRead)
	}
}
