package adplatform

import (
	"context"
)

type Campaign struct {
	ID          string
	Name        string
	Status      string
	DailyBudget float64
	Metrics     map[string]float64
}

type AdGroup struct {
	ID          string
	CampaignID  string
	Name        string
	Status      string
}

type Ad struct {
	ID         string
	AdGroupID  string
	Name       string
	Status     string
}

type PlatformConfig struct {
	AccessToken string
	AccountID  string
	CustomerID string
}

type PlatformAdapter interface {
	Connect(ctx context.Context, config PlatformConfig) error
	TestConnection(ctx context.Context, config PlatformConfig) error
	ListCampaigns(ctx context.Context, config PlatformConfig) ([]Campaign, error)
	GetCampaign(ctx context.Context, config PlatformConfig, campaignID string) (*Campaign, error)
	UpdateCampaignStatus(ctx context.Context, config PlatformConfig, campaignID, status string) error
	UpdateCampaignBudget(ctx context.Context, config PlatformConfig, campaignID string, budgetCents int64) error
	ListAdGroups(ctx context.Context, config PlatformConfig, campaignID string) ([]AdGroup, error)
	ListAds(ctx context.Context, config PlatformConfig, adGroupID string) ([]Ad, error)
	GetDailyMetrics(ctx context.Context, config PlatformConfig, startDate, endDate string) ([]map[string]interface{}, error)
}
