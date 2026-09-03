package handlers

import (
	"net/http"
	"time"
	"orbit/apps/api/internal/auth"
	"orbit/apps/api/internal/model"
	"orbit/apps/api/pkg/database"
	"orbit/apps/api/internal/tenant"
	"orbit/apps/api/internal/httputil"

	"github.com/gin-gonic/gin"
)

func ListConversations(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var conversations []model.Conversation
	if err := database.DB.Where("user_id = ? AND status = 'active'", userID).Order("created_at DESC").Limit(50).Find(&conversations).Error; err != nil {
		httputil.InternalError(c, "failed to list conversations")
		return
	}
	httputil.Success(c, gin.H{"items": conversations})
}

func GetConversationMessages(c *gin.Context) {
	userID, _ := c.Get("user_id")
	tenantID, _ := c.Get("tenant_id")
	id := c.Param("id")

	var conversation model.Conversation
	if err := database.DB.Where("id = ? AND user_id = ? AND tenant_id = ?", id, userID, tenantID).First(&conversation).Error; err != nil {
		httputil.NotFound(c, "conversation not found")
		return
	}

	var messages []model.Message
	if err := database.DB.Where("conversation_id = ?", id).Order("created_at ASC").Find(&messages).Error; err != nil {
		httputil.InternalError(c, "failed to fetch messages")
		return
	}
	httputil.Success(c, gin.H{"items": messages})
}

func ListPendingActions(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")
	var actions []model.AgentAction
	if err := database.DB.Where("tenant_id = ? AND status = 'pending'", tenantID).Order("created_at DESC").Find(&actions).Error; err != nil {
		httputil.InternalError(c, "failed to list actions")
		return
	}
	httputil.Success(c, gin.H{"items": actions})
}

func ApproveAction(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")
	id := c.Param("id")
	var action model.AgentAction
	if err := database.DB.Where("id = ? AND tenant_id = ? AND status = 'pending'", id, tenantID).First(&action).Error; err != nil {
		httputil.NotFound(c, "action not found or already processed")
		return
	}
	action.Status = "approved"
	now := time.Now()
	action.DecidedAt = &now
	if err := database.DB.Save(&action).Error; err != nil {
		httputil.InternalError(c, "failed to approve action")
		return
	}
	httputil.Success(c, gin.H{"message": "action approved", "action_id": action.ID})
}

func RejectAction(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")
	id := c.Param("id")
	var action model.AgentAction
	if err := database.DB.Where("id = ? AND tenant_id = ? AND status = 'pending'", id, tenantID).First(&action).Error; err != nil {
		httputil.NotFound(c, "action not found or already processed")
		return
	}
	action.Status = "rejected"
	now := time.Now()
	action.DecidedAt = &now
	if err := database.DB.Save(&action).Error; err != nil {
		httputil.InternalError(c, "failed to reject action")
		return
	}
	httputil.Success(c, gin.H{"message": "action rejected", "action_id": action.ID})
}

func RegisterAgentRoutes(r *gin.RouterGroup) {
	agent := r.Group("/agent")
	agent.Use(auth.NewAuthMiddleware().Handle())
	agent.Use(tenant.NewTenantMiddleware().Handle())
	{
		agent.POST("/chat", auth.RequirePermission("agent", "chat"), func(c *gin.Context) {
			c.Header("Content-Type", "text/event-stream")
			c.Header("Cache-Control", "no-cache")
			c.Header("Connection", "keep-alive")
			c.Status(http.StatusOK)
			c.Writer.WriteString("data: {\"message_id\":\"1\",\"delta\":\"????? Orbit ?????\"}\n\n")
			c.Writer.WriteString("data: {\"done\":true}\n\n")
			c.Writer.Flush()
		})
		agent.GET("/conversations", auth.RequirePermission("agent", "chat"), ListConversations)
		agent.GET("/conversations/:id/messages", auth.RequirePermission("agent", "chat"), GetConversationMessages)
		agent.GET("/actions", auth.RequirePermission("agent", "approve_action"), ListPendingActions)
		agent.POST("/actions/:id/approve", auth.RequirePermission("agent", "approve_action"), ApproveAction)
		agent.POST("/actions/:id/reject", auth.RequirePermission("agent", "reject_action"), RejectAction)
	}
}
