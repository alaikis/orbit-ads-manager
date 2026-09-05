package report

import (
	"context"
	"fmt"
	"strings"
	"orbit/apps/api/internal/model"
	"orbit/apps/api/pkg/database"
)

type ReportService struct{}

func NewReportService() *ReportService {
	return &ReportService{}
}

func (s *ReportService) GetSummary(ctx context.Context, tenantID uint64, scopeType string, scopeID *uint64, dateRange string, workspaceID *uint64) (map[string]interface{}, error) {
	var stats []model.DailyStat
	query := database.DB.Where("tenant_id = ? AND deleted_at IS NULL", tenantID)
	if workspaceID != nil {
		query = query.Where("workspace_id = ?", *workspaceID)
	} else {
		query = query.Where("workspace_id IS NULL")
	}
	if scopeID != nil {
		query = query.Where("scope_type = ? AND scope_id = ?", scopeType, *scopeID)
	}
	if err := query.Find(&stats).Error; err != nil {
		return nil, err
	}

	var totalSpend, totalImpressions, totalClicks, totalConversions float64
	for _, stat := range stats {
		if m, ok := stat.Metrics["spend"].(float64); ok {
			totalSpend += m
		}
		if m, ok := stat.Metrics["impressions"].(float64); ok {
			totalImpressions += m
		}
		if m, ok := stat.Metrics["clicks"].(float64); ok {
			totalClicks += m
		}
		if m, ok := stat.Metrics["conversions"].(float64); ok {
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

	return map[string]interface{}{
		"spend":       totalSpend,
		"impressions": totalImpressions,
		"clicks":      totalClicks,
		"ctr":         ctr,
		"cpc":         cpc,
		"conversions": totalConversions,
		"cpa":         cpa,
		"date_range":  dateRange,
	}, nil
}

func (s *ReportService) ExportCSV(ctx context.Context, tenantID uint64, scopeType string, scopeID *uint64, workspaceID *uint64) (string, error) {
	var stats []model.DailyStat
	query := database.DB.Where("tenant_id = ? AND deleted_at IS NULL", tenantID)
	if workspaceID != nil {
		query = query.Where("workspace_id = ?", *workspaceID)
	} else {
		query = query.Where("workspace_id IS NULL")
	}
	if scopeID != nil {
		query = query.Where("scope_type = ? AND scope_id = ?", scopeType, *scopeID)
	}
	if err := query.Find(&stats).Error; err != nil {
		return "", err
	}

	var sb strings.Builder
	sb.WriteString("date,platform,spend,impressions,clicks,ctr,cpc,conversions,cpa,roas\n")
	for _, s := range stats {
		m := s.Metrics
		sb.WriteString(fmt.Sprintf("%s,%s,%.2f,%.0f,%.0f,%.2f,%.2f,%.0f,%.2f,%.2f\n",
			s.Date, s.Platform,
			m["spend"], m["impressions"], m["clicks"],
			m["ctr"], m["cpc"], m["conversions"], m["cpa"], m["roas"]))
	}
	return sb.String(), nil
}

func (s *ReportService) GetTimeseries(ctx context.Context, tenantID uint64, scopeType string, step string, workspaceID *uint64) ([]map[string]interface{}, error) {
	var stats []model.DailyStat
	query := database.DB.Where("tenant_id = ? AND scope_type = ? AND deleted_at IS NULL", tenantID, scopeType)
	if workspaceID != nil {
		query = query.Where("workspace_id = ?", *workspaceID)
	} else {
		query = query.Where("workspace_id IS NULL")
	}
	if err := query.Order("date ASC").Limit(90).Find(&stats).Error; err != nil {
		return nil, err
	}

	points := make([]map[string]interface{}, len(stats))
	for i, s := range stats {
		points[i] = map[string]interface{}{"date": s.Date, "metrics": s.Metrics}
	}
	return points, nil
}
