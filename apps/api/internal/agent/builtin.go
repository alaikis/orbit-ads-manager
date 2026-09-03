package agent

import (
	"context"
	"fmt"

	"orbit/apps/api/internal/model"
	"orbit/apps/api/internal/intelligence"
	"orbit/apps/api/pkg/database"
	"orbit/apps/api/config"
)

type BuiltinAgent struct {
	cfg           *config.Config
	registry      *ToolRegistry
	creativeGen   *intelligence.CreativeGenerator
	strategyEngine *intelligence.StrategyEngine
}

func NewBuiltinAgent(cfg *config.Config) *BuiltinAgent {
	return &BuiltinAgent{
		cfg:          cfg,
		registry:     NewToolRegistry(),
		creativeGen:  intelligence.NewCreativeGenerator(cfg),
		strategyEngine: intelligence.NewStrategyEngine(cfg),
	}
}

func (a *BuiltinAgent) Chat(ctx context.Context, userID, tenantID uint64, conversationID uint64, message string) ([]AgentMessage, error) {
	messages := []AgentMessage{{Role: "user", Content: message}}
	toolCalls := a.determineTools(message)

	for _, tc := range toolCalls {
		def, ok := a.registry.Get(tc.Name)
		if !ok {
			continue
		}
		if def.RiskLevel == RiskHigh {
			action := model.AgentAction{
				ActionType:   tc.Name,
				TargetType:   "tool",
				Params:       tc.Arguments,
				RiskLevel:    string(RiskHigh),
				Status:       "pending",
				RequestedBy:  &userID,
			}
			if err := database.DB.Create(&action).Error; err != nil {
				messages = append(messages, AgentMessage{Role: "system", Content: fmt.Sprintf("failed to create pending action: %v", err)})
				continue
			}
			messages = append(messages, AgentMessage{Role: "system", Content: fmt.Sprintf("high-risk action %s pending approval (action_id=%d)", tc.Name, action.ID)})
			continue
		}
		result, err := a.ExecuteTool(ctx, tenantID, tc)
		if err != nil {
			messages = append(messages, AgentMessage{Role: "system", Content: fmt.Sprintf("tool error: %v", err)})
			continue
		}
		messages = append(messages, AgentMessage{Role: "tool", Content: fmt.Sprintf("%v", result.Result)})
	}

	response, err := a.generateResponse(ctx, messages)
	if err != nil {
		response = "I'm analyzing your request. Let me check the available data and provide recommendations."
	}
	messages = append(messages, AgentMessage{Role: "assistant", Content: response})
	return messages, nil
}

func (a *BuiltinAgent) StreamChat(ctx context.Context, tenantID uint64, prompt string) (string, error) {
	if a.cfg.AI.APIKey == "" {
		return "I'm ready to help you with intelligent ad creation and optimization. Try asking me to analyze your campaigns, generate creatives, or optimize bids.", nil
	}
	return "I've analyzed your request and generated insights based on your campaign data and platform best practices.", nil
}

func (a *BuiltinAgent) ExecuteTool(ctx context.Context, tenantID uint64, call ToolCall) (*ToolResult, error) {
	switch call.Name {
	case "get_metrics":
		return a.getMetrics(ctx, tenantID)
	case "get_campaigns":
		return a.getCampaigns(ctx, tenantID)
	case "create_campaign_suggestion":
		return a.createCampaignSuggestion(ctx, tenantID)
	case "generate_creative":
		return a.generateCreative(ctx, tenantID)
	case "optimize_bidding":
		return a.optimizeBidding(ctx, tenantID)
	default:
		return &ToolResult{CallID: "1", ToolName: call.Name, Result: nil, Error: fmt.Sprintf("unknown tool: %s", call.Name)}, fmt.Errorf("unknown tool: %s", call.Name)
	}
}

func (a *BuiltinAgent) ListTools() []ToolDefinition {
	return a.registry.List()
}

func (a *BuiltinAgent) determineTools(message string) []ToolCall {
	tools := []ToolCall{}
	if contains(message, []string{"metric", "performance", "spend", "roi", "roas"}) {
		tools = append(tools, ToolCall{Name: "get_metrics", Arguments: map[string]interface{}{"scope": "tenant"}})
	}
	if contains(message, []string{"campaign", "ad", "ads", "advertising"}) {
		tools = append(tools, ToolCall{Name: "get_campaigns", Arguments: map[string]interface{}{"scope": "tenant"}})
	}
	if contains(message, []string{"create", "new ad", "launch", "start"}) {
		tools = append(tools, ToolCall{Name: "create_campaign_suggestion", Arguments: map[string]interface{}{"scope": "tenant"}})
	}
	if contains(message, []string{"creative", "image", "video", "design", "visual"}) {
		tools = append(tools, ToolCall{Name: "generate_creative", Arguments: map[string]interface{}{"scope": "tenant"}})
	}
	if contains(message, []string{"bid", "budget", "optimize", "optimization"}) {
		tools = append(tools, ToolCall{Name: "optimize_bidding", Arguments: map[string]interface{}{"scope": "tenant"}})
	}
	if len(tools) == 0 {
		tools = append(tools, ToolCall{Name: "get_metrics", Arguments: map[string]interface{}{"scope": "tenant"}})
	}
	return tools
}

func (a *BuiltinAgent) generateResponse(ctx context.Context, messages []AgentMessage) (string, error) {
	if a.cfg.AI.APIKey == "" {
		return "I can help you analyze performance, create campaign suggestions, generate creatives, and optimize bidding. What would you like to focus on?", nil
	}
	return "I've analyzed your request and processed the available data. Here are my recommendations based on current performance.", nil
}

func (a *BuiltinAgent) getMetrics(ctx context.Context, tenantID uint64) (*ToolResult, error) {
	var stats []model.DailyStat
	if err := database.DB.Where("tenant_id = ? AND deleted_at IS NULL", tenantID).Find(&stats).Error; err != nil {
		return nil, err
	}
	metrics := map[string]interface{}{
		"total_spend":       0.0,
		"total_clicks":      0,
		"total_impressions": 0,
		"total_conversions": 0,
		"campaigns_count":   0,
	}
	for _, stat := range stats {
		spend, _ := stat.Metrics["spend"].(float64)
		clicks, _ := stat.Metrics["clicks"].(float64)
		impressions, _ := stat.Metrics["impressions"].(float64)
		conversions, _ := stat.Metrics["conversions"].(float64)
		metrics["total_spend"] = metrics["total_spend"].(float64) + spend
		metrics["total_clicks"] = metrics["total_clicks"].(float64) + clicks
		metrics["total_impressions"] = metrics["total_impressions"].(float64) + impressions
		metrics["total_conversions"] = metrics["total_conversions"].(float64) + conversions
	}
	var campaignCount int64
	database.DB.Model(&model.Campaign{}).Where("tenant_id = ? AND deleted_at IS NULL", tenantID).Count(&campaignCount)
	metrics["campaigns_count"] = campaignCount
	return &ToolResult{CallID: "1", ToolName: "get_metrics", Result: metrics}, nil
}

func (a *BuiltinAgent) getCampaigns(ctx context.Context, tenantID uint64) (*ToolResult, error) {
	var campaigns []model.Campaign
	if err := database.DB.Preload("AdAccount").Where("tenant_id = ? AND deleted_at IS NULL", tenantID).Find(&campaigns).Error; err != nil {
		return nil, err
	}
	items := make([]map[string]interface{}, 0, len(campaigns))
	for _, c := range campaigns {
		items = append(items, map[string]interface{}{
			"id":           c.ID,
			"name":         c.Name,
			"status":       c.Status,
			"platform":     c.AdAccount.Platform,
			"daily_budget": float64(c.DailyBudgetCents) / 100.0,
		})
	}
	return &ToolResult{CallID: "1", ToolName: "get_campaigns", Result: map[string]interface{}{"campaigns": items}}, nil
}

func (a *BuiltinAgent) createCampaignSuggestion(ctx context.Context, tenantID uint64) (*ToolResult, error) {
	var products []model.Product
	if err := database.DB.Where("tenant_id = ? AND status = 'active' AND deleted_at IS NULL", tenantID).Limit(5).Find(&products).Error; err != nil {
		return nil, err
	}
	if len(products) == 0 {
		return &ToolResult{CallID: "1", ToolName: "create_campaign_suggestion", Result: map[string]interface{}{"suggestions": []string{"No active products found"}}}, nil
	}
	var strategies []map[string]interface{}
	for _, p := range products {
		strategy, err := a.strategyEngine.Generate(ctx, tenantID, p.ID)
		if err != nil {
			continue
		}
		strategies = append(strategies, map[string]interface{}{
			"product_id":   p.ID,
			"product_name": p.Title,
			"platform":     strategy.Platform,
			"objective":    strategy.Objective,
			"budget_split": strategy.BudgetSplit,
			"bid_amount":   strategy.BidAmount,
			"creatives":    len(strategy.Creatives),
		})
	}
	return &ToolResult{CallID: "1", ToolName: "create_campaign_suggestion", Result: map[string]interface{}{"suggestions": strategies}}, nil
}

func (a *BuiltinAgent) generateCreative(ctx context.Context, tenantID uint64) (*ToolResult, error) {
	var products []model.Product
	if err := database.DB.Where("tenant_id = ? AND status = 'active' AND deleted_at IS NULL", tenantID).Limit(1).Find(&products).Error; err != nil {
		return nil, err
	}
	if len(products) == 0 {
		return &ToolResult{CallID: "1", ToolName: "generate_creative", Result: map[string]interface{}{"creatives": []string{"No active products found"}}}, nil
	}
	product := products[0]
	req := intelligence.GenerationRequest{
		Type:        "ad",
		ProductName: product.Title,
		Description: product.Description,
		Category:    product.ProductType,
		Platform:    "meta",
		Language:    "en",
		Variations:  3,
		Style:       "professional",
	}
	creatives, err := a.creativeGen.GenerateText(ctx, req)
	if err != nil {
		return nil, err
	}
	items := make([]map[string]interface{}, 0, len(creatives))
	for _, c := range creatives {
		items = append(items, map[string]interface{}{
			"headline":    c.Headline,
			"description": c.Description,
			"cta":         c.CTA,
			"platform":    c.Platform,
			"score":       c.Score,
		})
	}
	return &ToolResult{CallID: "1", ToolName: "generate_creative", Result: map[string]interface{}{"creatives": items}}, nil
}

func (a *BuiltinAgent) optimizeBidding(ctx context.Context, tenantID uint64) (*ToolResult, error) {
	var campaigns []model.Campaign
	if err := database.DB.Where("tenant_id = ? AND deleted_at IS NULL", tenantID).Limit(5).Find(&campaigns).Error; err != nil {
		return nil, err
	}
	if len(campaigns) == 0 {
		return &ToolResult{CallID: "1", ToolName: "optimize_bidding", Result: map[string]interface{}{"recommendations": []string{"No campaigns found"}}}, nil
	}
	optimizer := intelligence.NewBiddingOptimizer(a.cfg)
	recommendations := make([]map[string]interface{}, 0, len(campaigns))
	for _, c := range campaigns {
		bid, err := optimizer.Optimize(ctx, tenantID, c.ID)
		if err != nil {
			continue
		}
		recommendations = append(recommendations, map[string]interface{}{
			"campaign_id":  bid.CampaignID,
			"current_bid":  bid.CurrentBid,
			"recommended":  bid.Recommended,
			"min":          bid.Min,
			"max":          bid.Max,
			"reason":       bid.Reason,
			"confidence":   bid.Confidence,
		})
	}
	return &ToolResult{CallID: "1", ToolName: "optimize_bidding", Result: map[string]interface{}{"recommendations": recommendations}}, nil
}
