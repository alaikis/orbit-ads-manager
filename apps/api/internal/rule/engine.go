package rule

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"orbit/apps/api/internal/adplatform"
	"orbit/apps/api/internal/crypto"
	"orbit/apps/api/internal/model"
	"orbit/apps/api/pkg/database"
	"orbit/apps/api/config"
)

type ConditionSpec struct {
	Metric      string  `json:"metric"`
	Operator    string  `json:"operator"`
	Value       float64 `json:"value"`
	MinSpend    float64 `json:"min_spend,omitempty"`
	WindowDays  int     `json:"window_days,omitempty"`
	WindowType  string  `json:"window_type,omitempty"`
}

type ActionSpec struct {
	Type       string                 `json:"type"`
	Target     string                 `json:"target"`
	NewValue   float64                `json:"new_value,omitempty"`
	Level      string                 `json:"level,omitempty"`
	Message    string                 `json:"message,omitempty"`
}

type RuleEngine struct{}

func NewRuleEngine() *RuleEngine {
	return &RuleEngine{}
}

func (e *RuleEngine) Evaluate(ctx context.Context, rule *model.Rule) (*model.RuleExecution, error) {
	execution := &model.RuleExecution{
		RuleID:      rule.ID,
		TriggeredBy: "manual",
		Status:      "success",
	}

	condition, err := parseCondition(rule.ConditionSpec)
	if err != nil {
		execution.Status = "failed"
		execution.ErrorMessage = err.Error()
		return execution, err
	}

	matched, err := e.fetchAndMatch(rule.TenantID, condition)
	if err != nil {
		execution.Status = "failed"
		execution.ErrorMessage = err.Error()
		return execution, err
	}

	execution.MatchedObjects = map[string]interface{}{"items": matched}
	if len(matched) == 0 {
		execution.Status = "none"
		return execution, nil
	}

	action, err := parseAction(rule.ActionSpec)
	if err != nil {
		execution.Status = "failed"
		execution.ErrorMessage = err.Error()
		return execution, err
	}

	if err := e.executeActions(ctx, rule, action, matched); err != nil {
		execution.Status = "partial"
		execution.ErrorMessage = err.Error()
		return execution, nil
	}

	execution.ExecutedActions = map[string]interface{}{"action": action.Type, "targets": matched}
	execution.Status = "success"
	return execution, nil
}

func parseCondition(spec map[string]interface{}) (*ConditionSpec, error) {
	b, _ := json.Marshal(spec)
	var c ConditionSpec
	if err := json.Unmarshal(b, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

func parseAction(spec map[string]interface{}) (*ActionSpec, error) {
	b, _ := json.Marshal(spec)
	var a ActionSpec
	if err := json.Unmarshal(b, &a); err != nil {
		return nil, err
	}
	return &a, nil
}

func (e *RuleEngine) fetchAndMatch(tenantID uint64, cond *ConditionSpec) ([]map[string]interface{}, error) {
	var stats []model.DailyStat
	if err := database.DB.Where("tenant_id = ? AND date >= ? AND deleted_at IS NULL", tenantID, "2024-01-01").Find(&stats).Error; err != nil {
		return nil, err
	}

	var matched []map[string]interface{}
	for _, s := range stats {
		val, ok := s.Metrics[cond.Metric]
		if !ok {
			continue
		}
		fval, _ := val.(float64)
		if compare(fval, cond.Value, cond.Operator) {
			matched = append(matched, map[string]interface{}{"scope_type": s.ScopeType, "scope_id": s.ScopeID, "metric": cond.Metric, "value": fval})
		}
	}
	return matched, nil
}

func compare(actual, threshold float64, op string) bool {
	switch op {
	case "gt":
		return actual > threshold
	case "gte":
		return actual >= threshold
	case "lt":
		return actual < threshold
	case "lte":
		return actual <= threshold
	case "eq":
		return math.Abs(actual-threshold) < 0.0001
	case "neq":
		return math.Abs(actual-threshold) >= 0.0001
	default:
		return false
	}
}

func (e *RuleEngine) executeActions(ctx context.Context, rule *model.Rule, action *ActionSpec, targets []map[string]interface{}) error {
	var lastErr error
	for _, target := range targets {
		scopeType, _ := target["scope_type"].(string)
		scopeID, _ := target["scope_id"].(float64)

		adapter, platformConfig, err := getAdapterForRule(ctx, rule.TenantID, scopeType, uint64(scopeID))
		if err != nil {
			lastErr = err
			continue
		}

		switch action.Type {
		case "pause":
			if err := adapter.UpdateCampaignStatus(ctx, platformConfig, action.Target, "paused"); err != nil {
				lastErr = err
			}
		case "resume":
			if err := adapter.UpdateCampaignStatus(ctx, platformConfig, action.Target, "enabled"); err != nil {
				lastErr = err
			}
		case "update_budget":
			if err := adapter.UpdateCampaignBudget(ctx, platformConfig, action.Target, int64(action.NewValue*100)); err != nil {
				lastErr = err
			}
		case "update_bid":
			if err := adapter.UpdateCampaignBudget(ctx, platformConfig, action.Target, int64(action.NewValue*100)); err != nil {
				lastErr = err
			}
		case "notify":
			_ = createNotification(ctx, rule.TenantID, action.Message)
		}
	}
	return lastErr
}

func getAdapterForRule(ctx context.Context, tenantID uint64, scopeType string, scopeID uint64) (adplatform.PlatformAdapter, adplatform.PlatformConfig, error) {
	var account model.AdAccount
	if err := database.DB.Preload("Conn").Where("tenant_id = ? AND id = ? AND deleted_at IS NULL", tenantID, scopeID).First(&account).Error; err != nil {
		var campaign model.Campaign
		if err := database.DB.Where("tenant_id = ? AND id = ? AND deleted_at IS NULL", tenantID, scopeID).First(&campaign).Error; err != nil {
			return nil, adplatform.PlatformConfig{}, fmt.Errorf("resource not found for scope %s:%d", scopeType, scopeID)
		}
		if err := database.DB.Preload("Conn").Where("id = ?", campaign.AdAccountID).First(&account).Error; err != nil {
			return nil, adplatform.PlatformConfig{}, fmt.Errorf("ad account not found for campaign %d", campaign.AdAccountID)
		}
	}

	var conn model.PlatformConn
	if err := database.DB.Where("id = ?", account.ConnID).First(&conn).Error; err != nil {
		return nil, adplatform.PlatformConfig{}, fmt.Errorf("platform connection not found: %w", err)
	}

	var oauthToken model.OAuthToken
	if err := database.DB.Where("conn_id = ?", conn.ID).First(&oauthToken).Error; err != nil {
		return nil, adplatform.PlatformConfig{}, fmt.Errorf("oauth token not found: %w", err)
	}

	accessToken, err := crypto.Decrypt(oauthToken.EncryptedToken1)
	if err != nil {
		return nil, adplatform.PlatformConfig{}, fmt.Errorf("failed to decrypt access token: %w", err)
	}

	cfg := config.AppCfg
	var adapter adplatform.PlatformAdapter
	switch conn.Platform {
	case "meta":
		adapter = adplatform.NewMetaAdapter(cfg)
	case "google":
		adapter = adplatform.NewGoogleAdapter(cfg)
	default:
		return nil, adplatform.PlatformConfig{}, fmt.Errorf("unsupported platform: %s", conn.Platform)
	}

	platformConfig := adplatform.PlatformConfig{
		AccessToken: accessToken,
		AccountID:   account.ExternalID,
		CustomerID:  account.ExternalID,
	}

	return adapter, platformConfig, nil
}

func createNotification(ctx context.Context, tenantID uint64, message string) error {
	if message == "" {
		return nil
	}

	notification := model.Notification{
		BaseModel: model.BaseModel{TenantID: tenantID},
		Type:      "rule_action",
		Title:     "Rule Executed",
		Body:      message,
	}
	return database.DB.Create(&notification).Error
}
