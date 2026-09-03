package adplatform

import (
	"context"
	"fmt"
	"orbit/apps/api/config"
)

type BingAdapter struct {
	config *config.Config
}

func NewBingAdapter(cfg *config.Config) *BingAdapter {
	return &BingAdapter{config: cfg}
}

func (a *BingAdapter) Connect(ctx context.Context, config interface{}) error {
	return nil
}

func (a *BingAdapter) TestConnection(ctx context.Context) error {
	return nil
}

func (a *BingAdapter) ListCampaigns(ctx context.Context, accountID string) ([]Campaign, error) {
	return nil, fmt.Errorf("not implemented")
}

func (a *BingAdapter) GetCampaign(ctx context.Context, accountID, campaignID string) (*Campaign, error) {
	return nil, fmt.Errorf("not implemented")
}

func (a *BingAdapter) UpdateCampaignStatus(ctx context.Context, accountID, campaignID, status string) error {
	return fmt.Errorf("not implemented")
}

func (a *BingAdapter) UpdateCampaignBudget(ctx context.Context, accountID, campaignID string, budgetCents int64) error {
	return fmt.Errorf("not implemented")
}

func (a *BingAdapter) ListAdGroups(ctx context.Context, accountID, campaignID string) ([]AdGroup, error) {
	return nil, fmt.Errorf("not implemented")
}

func (a *BingAdapter) ListAds(ctx context.Context, accountID, adGroupID string) ([]Ad, error) {
	return nil, fmt.Errorf("not implemented")
}

func (a *BingAdapter) GetDailyMetrics(ctx context.Context, accountID, startDate, endDate string) ([]map[string]interface{}, error) {
	return nil, fmt.Errorf("not implemented")
}
