package handlers

import (
	"strconv"

	"orbit/apps/api/internal/auth"
	"orbit/apps/api/internal/model"
	"orbit/apps/api/internal/report"
	"orbit/apps/api/pkg/database"
	"orbit/apps/api/internal/tenant"
	"orbit/apps/api/internal/httputil"

	"github.com/gin-gonic/gin"
)

func GetReportsSummary(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")
	workspaceID := c.Query("workspace_id")
	scopeType := c.DefaultQuery("scope_type", "")
	scopeIDStr := c.Query("scope_id")

	var wsID *uint64
	if workspaceID != "" {
		id, _ := strconv.ParseUint(workspaceID, 10, 64)
		wsID = &id
	}
	var sID *uint64
	if scopeIDStr != "" {
		id, _ := strconv.ParseUint(scopeIDStr, 10, 64)
		sID = &id
	}

	svc := report.NewReportService()
	summary, err := svc.GetSummary(c.Request.Context(), tenantID.(uint64), scopeType, sID, "last_30d", wsID)
	if err != nil {
		httputil.InternalError(c, "failed to get summary")
		return
	}
	httputil.Success(c, summary)
}

func ListReports(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")
	var schedules []model.ReportSchedule
	query := database.DB.Where("tenant_id = ? AND deleted_at IS NULL", tenantID)
	if err := query.Order("created_at DESC").Find(&schedules).Error; err != nil {
		httputil.InternalError(c, "failed to list reports")
		return
	}
	httputil.Success(c, gin.H{"items": schedules})
}

func ExportReport(c *gin.Context) {
	httputil.Accepted(c, gin.H{"job_id": "export-123", "status": "queued"})
}

func ListSchedules(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")
	workspaceID := c.Query("workspace_id")
	var schedules []model.ReportSchedule
	query := database.DB.Where("tenant_id = ? AND deleted_at IS NULL", tenantID)
	if workspaceID != "" {
		query = query.Where("workspace_id = ?", workspaceID)
	} else {
		query = query.Where("workspace_id IS NULL")
	}
	if err := query.Order("created_at DESC").Find(&schedules).Error; err != nil {
		httputil.InternalError(c, "failed to list schedules")
		return
	}
	httputil.Success(c, gin.H{"items": schedules})
}

func CreateSchedule(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")
	var schedule model.ReportSchedule
	if err := c.ShouldBindJSON(&schedule); err != nil {
		httputil.BadRequest(c, err.Error(), nil)
		return
	}
	schedule.TenantID = tenantID.(uint64)
	if schedule.Name == "" {
		httputil.BadRequest(c, "name is required", nil)
		return
	}
	if err := database.DB.Create(&schedule).Error; err != nil {
		httputil.InternalError(c, "failed to create schedule")
		return
	}
	httputil.Created(c, schedule)
}

func RegisterReportRoutes(r *gin.RouterGroup) {
	reports := r.Group("/reports")
	reports.Use(auth.NewAuthMiddleware().Handle())
	reports.Use(tenant.NewTenantMiddleware().Handle())
	{
		reports.GET("/summary", auth.RequirePermission("report", "read"), GetReportsSummary)
		reports.GET("", auth.RequirePermission("report", "read"), ListReports)
		reports.POST("/export", auth.RequirePermission("report", "export"), ExportReport)
		reports.GET("/schedules", auth.RequirePermission("report", "read"), ListSchedules)
		reports.POST("/schedules", auth.RequirePermission("report", "schedule"), CreateSchedule)
	}
}
