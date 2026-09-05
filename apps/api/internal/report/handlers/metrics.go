package handlers

import (
	"orbit/apps/api/internal/auth"
	"orbit/apps/api/internal/model"
	"orbit/apps/api/pkg/database"
	"orbit/apps/api/internal/tenant"
	"orbit/apps/api/internal/httputil"

	"github.com/gin-gonic/gin"
)

func GetMetricsSummary(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")
	scopeType := c.DefaultQuery("scope_type", "tenant")
	scopeID := c.Query("scope_id")
	workspaceID := c.Query("workspace_id")
	_ = c.DefaultQuery("date_range", "7d")

	var stats []model.DailyStat
	query := database.DB.Where("tenant_id = ? AND deleted_at IS NULL", tenantID)
	if workspaceID != "" {
		query = query.Where("workspace_id = ?", workspaceID)
	} else {
		query = query.Where("workspace_id IS NULL")
	}
	if scopeID != "" {
		query = query.Where("scope_type = ? AND scope_id = ?", scopeType, scopeID)
	}
	query = query.Order("date DESC").Limit(30)

	if err := query.Find(&stats).Error; err != nil {
		httputil.InternalError(c, "failed to fetch metrics")
		return
	}

	aggregated := aggregateStats(stats)
	httputil.Success(c, aggregated)
}

func GetTimeseries(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")
	scopeType := c.DefaultQuery("scope", "tenant")
	step := c.DefaultQuery("step", "day")
	workspaceID := c.Query("workspace_id")

	var stats []model.DailyStat
	query := database.DB.Where("tenant_id = ? AND scope_type = ? AND deleted_at IS NULL", tenantID, scopeType)
	if workspaceID != "" {
		query = query.Where("workspace_id = ?", workspaceID)
	} else {
		query = query.Where("workspace_id IS NULL")
	}
	query = query.Order("date ASC").Limit(90)

	if err := query.Find(&stats).Error; err != nil {
		httputil.InternalError(c, "failed to fetch timeseries")
		return
	}

	points := make([]gin.H, len(stats))
	for i, s := range stats {
		points[i] = gin.H{"date": s.Date, "metrics": s.Metrics}
	}
	httputil.Success(c, gin.H{"points": points, "step": step})
}

func aggregateStats(stats []model.DailyStat) gin.H {
	var totalSpend, totalImpressions, totalClicks, totalConversions float64
	for _, s := range stats {
		if m, ok := s.Metrics["spend"].(float64); ok {
			totalSpend += m
		}
		if m, ok := s.Metrics["impressions"].(float64); ok {
			totalImpressions += m
		}
		if m, ok := s.Metrics["clicks"].(float64); ok {
			totalClicks += m
		}
		if m, ok := s.Metrics["conversions"].(float64); ok {
			totalConversions += m
		}
	}
	ctr := 0.0
	if totalImpressions > 0 {
		ctr = (totalClicks / totalImpressions) * 100
	}
	cpc := 0.0
	if totalClicks > 0 {
		cpc = totalSpend / totalClicks
	}
	cpa := 0.0
	if totalConversions > 0 {
		cpa = totalSpend / totalConversions
	}
	roas := 0.0
	if totalSpend > 0 {
		if m, ok := stats[0].Metrics["conversion_value"].(float64); ok {
			roas = m / totalSpend
		}
	}
	return gin.H{
		"spend": totalSpend,
		"impressions": totalImpressions,
		"clicks": totalClicks,
		"ctr": ctr,
		"cpc": cpc,
		"conversions": totalConversions,
		"cpa": cpa,
		"roas": roas,
	}
}

func RegisterMetricsRoutes(r *gin.RouterGroup) {
	metrics := r.Group("/metrics")
	metrics.Use(auth.NewAuthMiddleware().Handle())
	metrics.Use(tenant.NewTenantMiddleware().Handle())
	{
		metrics.GET("/summary", auth.RequirePermission("report", "read"), GetMetricsSummary)
		metrics.GET("/timeseries", auth.RequirePermission("report", "read"), GetTimeseries)
	}
}
