package provider

import (
	"time"
	"orbit/apps/api/pkg/database"
)

type UsageRecord struct {
	ProviderID   uint64
	WorkspaceID  uint64
	TenantID     uint64
	TokensUsed   int64
	RequestsCount int
	ErrorsCount  int
	CostCents    int64
}

type QuotaLimit struct {
	MaxTokensPerMonth int64
	MaxRequestsPerDay int
	MaxCostCentsPerMonth int64
}

var DefaultQuota = QuotaLimit{
	MaxTokensPerMonth:    1000000,
	MaxRequestsPerDay:    10000,
	MaxCostCentsPerMonth: 50000,
}

func RecordUsage(record UsageRecord) error {
	now := time.Now()
	periodStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	periodEnd := periodStart.AddDate(0, 1, 0)

	var usage ProviderUsage
	if err := database.DB.Where("provider_id = ? AND period_start = ? AND period_end = ?", record.ProviderID, periodStart, periodEnd).First(&usage).Error; err != nil {
		usage = ProviderUsage{
			ProviderID:   record.ProviderID,
			WorkspaceID:  record.WorkspaceID,
			TenantID:     &record.TenantID,
			PeriodStart:  periodStart,
			PeriodEnd:    periodEnd,
			TokensUsed:   record.TokensUsed,
			RequestsCount: record.RequestsCount,
			ErrorsCount:  record.ErrorsCount,
			CostCents:    record.CostCents,
		}
		return database.DB.Create(&usage).Error
	}

	usage.TokensUsed += record.TokensUsed
	usage.RequestsCount += record.RequestsCount
	usage.ErrorsCount += record.ErrorsCount
	usage.CostCents += record.CostCents
	usage.UpdatedAt = now
	return database.DB.Save(&usage).Error
}

func CheckQuota(workspaceID uint64, limit QuotaLimit) (bool, string) {
	now := time.Now()
	periodStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	var totalTokens int64
	var totalCost int64
	database.DB.Model(&ProviderUsage{}).Where("workspace_id = ? AND period_start >= ?", workspaceID, periodStart).Row().Scan(&totalTokens, &totalCost)

	if totalTokens >= limit.MaxTokensPerMonth {
		return false, "monthly token quota exceeded"
	}
	if totalCost >= limit.MaxCostCentsPerMonth {
		return false, "monthly cost quota exceeded"
	}
	return true, ""
}
