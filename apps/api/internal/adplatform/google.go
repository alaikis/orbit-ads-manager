package adplatform

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"orbit/apps/api/config"
)

type GoogleAdapter struct {
	config *config.Config
}

func NewGoogleAdapter(cfg *config.Config) *GoogleAdapter {
	return &GoogleAdapter{config: cfg}
}

func (a *GoogleAdapter) Connect(ctx context.Context, cfg PlatformConfig) error {
	return a.TestConnection(ctx, cfg)
}

func (a *GoogleAdapter) TestConnection(ctx context.Context, cfg PlatformConfig) error {
	url := fmt.Sprintf("https://googleads.googleapis.com/v14/customers/%s:listAccessibleCustomers", cfg.CustomerID)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+cfg.AccessToken)
	req.Header.Set("developer-token", a.config.OAuth.GoogleDeveloperToken)

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("google ads API returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

func (a *GoogleAdapter) ListCampaigns(ctx context.Context, cfg PlatformConfig) ([]Campaign, error) {
	customerID := strings.ReplaceAll(cfg.CustomerID, "-", "")
	query := "SELECT campaign.id, campaign.name, campaign.status, campaign.daily_budget.micros FROM campaign WHERE campaign.status != 'REMOVED'"

	result, err := a.search(ctx, cfg, customerID, query)
	if err != nil {
		return nil, err
	}

	campaigns := make([]Campaign, 0, len(result))
	for _, row := range result {
		campaign := Campaign{
			ID:     getString(row, "campaign.id"),
			Name:   getString(row, "campaign.name"),
			Status: getString(row, "campaign.status"),
		}

		if micros := getString(row, "campaign.daily_budget.micros"); micros != "" {
			var budgetMicros int64
			_, _ = fmt.Sscanf(micros, "%d", &budgetMicros)
			campaign.DailyBudget = float64(budgetMicros) / 1_000_000.0
		}

		campaigns = append(campaigns, campaign)
	}

	return campaigns, nil
}

func (a *GoogleAdapter) GetCampaign(ctx context.Context, cfg PlatformConfig, campaignID string) (*Campaign, error) {
	customerID := strings.ReplaceAll(cfg.CustomerID, "-", "")
	query := fmt.Sprintf("SELECT campaign.id, campaign.name, campaign.status, campaign.daily_budget.micros FROM campaign WHERE campaign.id = '%s'", campaignID)

	result, err := a.search(ctx, cfg, customerID, query)
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("campaign not found")
	}

	row := result[0]
	campaign := &Campaign{
		ID:     getString(row, "campaign.id"),
		Name:   getString(row, "campaign.name"),
		Status: getString(row, "campaign.status"),
	}

	if micros := getString(row, "campaign.daily_budget.micros"); micros != "" {
		var budgetMicros int64
		_, _ = fmt.Sscanf(micros, "%d", &budgetMicros)
		campaign.DailyBudget = float64(budgetMicros) / 1_000_000.0
	}

	return campaign, nil
}

func (a *GoogleAdapter) UpdateCampaignStatus(ctx context.Context, cfg PlatformConfig, campaignID, status string) error {
	customerID := strings.ReplaceAll(cfg.CustomerID, "-", "")

	mutations := map[string]interface{}{
		"operations": []map[string]interface{}{
			{
				"update": map[string]interface{}{
					"resourceName": fmt.Sprintf("customers/%s/campaigns/%s", customerID, campaignID),
					"status":       strings.ToUpper(status),
				},
				"updateMask": "status",
			},
		},
	}

	return a.mutate(ctx, cfg, customerID, mutations)
}

func (a *GoogleAdapter) UpdateCampaignBudget(ctx context.Context, cfg PlatformConfig, campaignID string, budgetCents int64) error {
	customerID := strings.ReplaceAll(cfg.CustomerID, "-", "")
	micros := budgetCents * 10000

	mutations := map[string]interface{}{
		"operations": []map[string]interface{}{
			{
				"update": map[string]interface{}{
					"resourceName": fmt.Sprintf("customers/%s/campaigns/%s", customerID, campaignID),
					"dailyBudget": map[string]interface{}{
						"micros": micros,
					},
				},
				"updateMask": "campaign.daily_budget",
			},
		},
	}

	return a.mutate(ctx, cfg, customerID, mutations)
}

func (a *GoogleAdapter) ListAdGroups(ctx context.Context, cfg PlatformConfig, campaignID string) ([]AdGroup, error) {
	customerID := strings.ReplaceAll(cfg.CustomerID, "-", "")
	query := fmt.Sprintf("SELECT ad_group.id, ad_group.name, ad_group.status, ad_group.campaign FROM ad_group WHERE ad_group.campaign = 'customers/%s/campaigns/%s' AND ad_group.status != 'REMOVED'", customerID, campaignID)

	result, err := a.search(ctx, cfg, customerID, query)
	if err != nil {
		return nil, err
	}

	adGroups := make([]AdGroup, 0, len(result))
	for _, row := range result {
		adGroups = append(adGroups, AdGroup{
			ID:         getString(row, "ad_group.id"),
			CampaignID: campaignID,
			Name:       getString(row, "ad_group.name"),
			Status:     getString(row, "ad_group.status"),
		})
	}

	return adGroups, nil
}

func (a *GoogleAdapter) ListAds(ctx context.Context, cfg PlatformConfig, adGroupID string) ([]Ad, error) {
	customerID := strings.ReplaceAll(cfg.CustomerID, "-", "")
	query := fmt.Sprintf("SELECT ad_group_ad.ad.id, ad_group_ad.ad.name, ad_group_ad.ad.status, ad_group_ad.ad_group FROM ad_group_ad WHERE ad_group_ad.ad_group = 'customers/%s/adGroups/%s'", customerID, adGroupID)

	result, err := a.search(ctx, cfg, customerID, query)
	if err != nil {
		return nil, err
	}

	ads := make([]Ad, 0, len(result))
	for _, row := range result {
		ads = append(ads, Ad{
			ID:        getString(row, "ad_group_ad.ad.id"),
			AdGroupID: adGroupID,
			Name:      getString(row, "ad_group_ad.ad.name"),
			Status:    getString(row, "ad_group_ad.ad.status"),
		})
	}

	return ads, nil
}

func (a *GoogleAdapter) GetDailyMetrics(ctx context.Context, cfg PlatformConfig, startDate, endDate string) ([]map[string]interface{}, error) {
	customerID := strings.ReplaceAll(cfg.CustomerID, "-", "")
	query := fmt.Sprintf("SELECT campaign.id, campaign.name, metrics.impressions, metrics.clicks, metrics.cost_micros, metrics.average_cpc_micros, metrics.ctr FROM campaign WHERE segments.date >= '%s' AND segments.date <= '%s'", startDate, endDate)

	result, err := a.search(ctx, cfg, customerID, query)
	if err != nil {
		return nil, err
	}

	metrics := make([]map[string]interface{}, 0, len(result))
	for _, row := range result {
		metric := map[string]interface{}{
			"campaign_id": getString(row, "campaign.id"),
			"name":        getString(row, "campaign.name"),
			"impressions": getString(row, "metrics.impressions"),
			"clicks":      getString(row, "metrics.clicks"),
		}

		if costMicros := getString(row, "metrics.cost_micros"); costMicros != "" {
			var cost int64
			_, _ = fmt.Sscanf(costMicros, "%d", &cost)
			metric["spend"] = float64(cost) / 1_000_000.0
		}

		if cpcMicros := getString(row, "metrics.average_cpc_micros"); cpcMicros != "" {
			var cpc int64
			_, _ = fmt.Sscanf(cpcMicros, "%d", &cpc)
			metric["cpc"] = float64(cpc) / 1_000_000.0
		}

		if ctr := getString(row, "metrics.ctr"); ctr != "" {
			metric["ctr"] = ctr
		}

		metrics = append(metrics, metric)
	}

	return metrics, nil
}

func (a *GoogleAdapter) search(ctx context.Context, cfg PlatformConfig, customerID, query string) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("https://googleads.googleapis.com/v14/customers/%s/googleAds:search", customerID)

	payload := map[string]interface{}{
		"query": query,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+cfg.AccessToken)
	req.Header.Set("developer-token", a.config.OAuth.GoogleDeveloperToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google ads API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Results []map[string]interface{} `json:"results"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, err
	}

	return result.Results, nil
}

func (a *GoogleAdapter) mutate(ctx context.Context, cfg PlatformConfig, customerID string, mutations map[string]interface{}) error {
	url := fmt.Sprintf("https://googleads.googleapis.com/v14/customers/%s/campaigns:mutate", customerID)

	body, err := json.Marshal(mutations)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(body)))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+cfg.AccessToken)
	req.Header.Set("developer-token", a.config.OAuth.GoogleDeveloperToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("google ads API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

func getString(m map[string]interface{}, path string) string {
	parts := strings.Split(path, ".")
	var current interface{} = m

	for _, part := range parts {
		if part == "" {
			continue
		}
		if mapVal, ok := current.(map[string]interface{}); ok {
			if val, exists := mapVal[part]; exists {
				current = val
			} else {
				return ""
			}
		} else {
			return ""
		}
	}

	if str, ok := current.(string); ok {
		return str
	}
	return ""
}
