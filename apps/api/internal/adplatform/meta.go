package adplatform

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"orbit/apps/api/config"
)

type MetaAdapter struct {
	config *config.Config
}

func NewMetaAdapter(cfg *config.Config) *MetaAdapter {
	return &MetaAdapter{config: cfg}
}

func (a *MetaAdapter) Connect(ctx context.Context, cfg PlatformConfig) error {
	return a.TestConnection(ctx, cfg)
}

func (a *MetaAdapter) TestConnection(ctx context.Context, cfg PlatformConfig) error {
	url := fmt.Sprintf("https://graph.facebook.com/v19.0/me?access_token=%s", cfg.AccessToken)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("meta API returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

func (a *MetaAdapter) ListCampaigns(ctx context.Context, cfg PlatformConfig) ([]Campaign, error) {
	accountID := cfg.AccountID
	if accountID == "" {
		return nil, fmt.Errorf("ad account ID is required")
	}

	url := fmt.Sprintf("https://graph.facebook.com/v19.0/act_%s/campaigns?fields=id,name,status,daily_budget,objective&limit=100&access_token=%s", accountID, cfg.AccessToken)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("meta API returned status %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data []struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			Status      string `json:"status"`
			DailyBudget string `json:"daily_budget"`
			Objective   string `json:"objective"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	campaigns := make([]Campaign, 0, len(result.Data))
	for _, c := range result.Data {
		budget := parseMetaBudget(c.DailyBudget)
		campaigns = append(campaigns, Campaign{
			ID:          c.ID,
			Name:        c.Name,
			Status:      c.Status,
			DailyBudget: budget,
			Metrics:     map[string]float64{"objective": 0},
		})
	}

	return campaigns, nil
}

func (a *MetaAdapter) GetCampaign(ctx context.Context, cfg PlatformConfig, campaignID string) (*Campaign, error) {
	url := fmt.Sprintf("https://graph.facebook.com/v19.0/%s?fields=id,name,status,daily_budget,objective&access_token=%s", campaignID, cfg.AccessToken)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("meta API returned status %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Status      string `json:"status"`
		DailyBudget string `json:"daily_budget"`
		Objective   string `json:"objective"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &Campaign{
		ID:          result.ID,
		Name:        result.Name,
		Status:      result.Status,
		DailyBudget: parseMetaBudget(result.DailyBudget),
		Metrics:     map[string]float64{"objective": 0},
	}, nil
}

func (a *MetaAdapter) UpdateCampaignStatus(ctx context.Context, cfg PlatformConfig, campaignID, status string) error {
	url := fmt.Sprintf("https://graph.facebook.com/v19.0/%s?access_token=%s", campaignID, cfg.AccessToken)
	payload := map[string]interface{}{
		"status": status,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("meta API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

func (a *MetaAdapter) UpdateCampaignBudget(ctx context.Context, cfg PlatformConfig, campaignID string, budgetCents int64) error {
	url := fmt.Sprintf("https://graph.facebook.com/v19.0/%s?access_token=%s", campaignID, cfg.AccessToken)
	budget := fmt.Sprintf("%d", budgetCents)
	payload := map[string]interface{}{
		"daily_budget": budget,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("meta API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

func (a *MetaAdapter) ListAdGroups(ctx context.Context, cfg PlatformConfig, campaignID string) ([]AdGroup, error) {
	url := fmt.Sprintf("https://graph.facebook.com/v19.0/%s/adsets?fields=id,name,status,campaign_id&limit=100&access_token=%s", campaignID, cfg.AccessToken)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("meta API returned status %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data []struct {
			ID         string `json:"id"`
			Name       string `json:"name"`
			Status     string `json:"status"`
			CampaignID string `json:"campaign_id"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	adGroups := make([]AdGroup, 0, len(result.Data))
	for _, ag := range result.Data {
		adGroups = append(adGroups, AdGroup{
			ID:         ag.ID,
			CampaignID: ag.CampaignID,
			Name:       ag.Name,
			Status:     ag.Status,
		})
	}

	return adGroups, nil
}

func (a *MetaAdapter) ListAds(ctx context.Context, cfg PlatformConfig, adGroupID string) ([]Ad, error) {
	url := fmt.Sprintf("https://graph.facebook.com/v19.0/%s/ads?fields=id,name,status,adset_id&limit=100&access_token=%s", adGroupID, cfg.AccessToken)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("meta API returned status %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data []struct {
			ID       string `json:"id"`
			Name     string `json:"name"`
			Status   string `json:"status"`
			AdsetID  string `json:"adset_id"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	ads := make([]Ad, 0, len(result.Data))
	for _, ad := range result.Data {
		ads = append(ads, Ad{
			ID:         ad.ID,
			AdGroupID:  ad.AdsetID,
			Name:       ad.Name,
			Status:     ad.Status,
		})
	}

	return ads, nil
}

func (a *MetaAdapter) GetDailyMetrics(ctx context.Context, cfg PlatformConfig, startDate, endDate string) ([]map[string]interface{}, error) {
	accountID := cfg.AccountID
	if accountID == "" {
		return nil, fmt.Errorf("ad account ID is required")
	}

	url := fmt.Sprintf("https://graph.facebook.com/v19.0/act_%s/insights?fields=spend,impressions,clicks,actions,cpc,cpm,ctr&time_range={\"since\":\"%s\",\"until\":\"%s\"}&access_token=%s", accountID, startDate, endDate, cfg.AccessToken)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("meta API returned status %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data []map[string]interface{} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

func parseMetaBudget(budget string) float64 {
	if budget == "" {
		return 0
	}
	var cents int64
	_, _ = fmt.Sscanf(budget, "%d", &cents)
	return float64(cents) / 100.0
}
