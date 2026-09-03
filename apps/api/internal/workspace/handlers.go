package workspace

import (
	"orbit/apps/api/internal/auth"
	"orbit/apps/api/pkg/database"
	"orbit/apps/api/internal/tenant"
	"orbit/apps/api/internal/httputil"

	"github.com/gin-gonic/gin"
)

func ListWorkspaces(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var items []Workspace
	if err := database.DB.Where("owner_user_id = ?", userID).Find(&items).Error; err != nil {
		httputil.InternalError(c, "failed to list workspaces")
		return
	}
	httputil.Success(c, gin.H{"items": items})
}

func CreateWorkspace(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var req struct {
		Name string `json:"name" binding:"required"`
		Plan string `json:"plan" binding:"required,oneof=beta pro enterprise"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, err.Error(), nil)
		return
	}

	ws := Workspace{
		Name:        req.Name,
		OwnerUserID: userID.(uint64),
		Plan:        req.Plan,
	}
	if err := database.DB.Create(&ws).Error; err != nil {
		httputil.InternalError(c, "failed to create workspace")
		return
	}
	httputil.Created(c, ws)
}

func GetWorkspace(c *gin.Context) {
	id := c.Param("id")
	var ws Workspace
	if err := database.DB.First(&ws, id).Error; err != nil {
		httputil.NotFound(c, "workspace not found")
		return
	}
	httputil.Success(c, ws)
}

func UpdateWorkspace(c *gin.Context) {
	id := c.Param("id")
	var ws Workspace
	if err := database.DB.First(&ws, id).Error; err != nil {
		httputil.NotFound(c, "workspace not found")
		return
	}
	var req struct {
		Name     string                 `json:"name"`
		Plan     string                 `json:"plan"`
		Settings map[string]interface{} `json:"settings"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, err.Error(), nil)
		return
	}
	if req.Name != "" {
		ws.Name = req.Name
	}
	if req.Plan != "" {
		ws.Plan = req.Plan
	}
	if req.Settings != nil {
		ws.Settings = req.Settings
	}
	if err := database.DB.Save(&ws).Error; err != nil {
		httputil.InternalError(c, "failed to update workspace")
		return
	}
	httputil.Success(c, ws)
}

func DeleteWorkspace(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Delete(&Workspace{}, id).Error; err != nil {
		httputil.InternalError(c, "failed to delete workspace")
		return
	}
	httputil.Success(c, gin.H{"message": "deleted"})
}

func SwitchWorkspace(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var req struct {
		WorkspaceID uint64 `json:"workspace_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, err.Error(), nil)
		return
	}
	var ws Workspace
	if err := database.DB.Where("id = ? AND owner_user_id = ?", req.WorkspaceID, userID).First(&ws).Error; err != nil {
		httputil.NotFound(c, "workspace not found")
		return
	}
	httputil.Success(c, gin.H{"message": "switched", "workspace_id": ws.ID})
}

func GetCurrentWorkspace(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var ws Workspace
	if err := database.DB.Where("owner_user_id = ?", userID).Order("created_at ASC").First(&ws).Error; err != nil {
		httputil.Success(c, nil)
		return
	}
	httputil.Success(c, ws)
}

func RegisterWorkspaceRoutes(r *gin.RouterGroup) {
	workspaces := r.Group("/workspaces")
	workspaces.Use(auth.NewAuthMiddleware().Handle())
	workspaces.Use(tenant.NewTenantMiddleware().Handle())
	{
		workspaces.GET("/current", GetCurrentWorkspace)
		workspaces.GET("", ListWorkspaces)
		workspaces.POST("", CreateWorkspace)
		workspaces.GET("/:id", GetWorkspace)
		workspaces.PATCH("/:id", UpdateWorkspace)
		workspaces.DELETE("/:id", DeleteWorkspace)
		workspaces.POST("/:id/switch", SwitchWorkspace)
	}
}
