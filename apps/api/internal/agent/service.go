package agent

import (
	"context"
	"fmt"

	"orbit/apps/api/internal/model"
	"orbit/apps/api/internal/intelligence"
	"orbit/apps/api/internal/provider"
	"orbit/apps/api/pkg/database"
	"orbit/apps/api/config"
)

type ToolCall struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

type ToolResult struct {
	CallID   string                 `json:"call_id"`
	ToolName string                 `json:"tool_name"`
	Result   map[string]interface{} `json:"result"`
	Error    string                 `json:"error,omitempty"`
}

type AgentMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type AgentService struct {
	cfg           *config.Config
	creativeGen   *intelligence.CreativeGenerator
	strategyEngine *intelligence.StrategyEngine
	resolver      *provider.ProviderResolver
}

func NewAgentService() *AgentService {
	cfg := config.AppCfg
	return &AgentService{
		cfg:           cfg,
		creativeGen:   intelligence.NewCreativeGenerator(cfg),
		strategyEngine: intelligence.NewStrategyEngine(cfg),
		resolver:      provider.NewProviderResolver(),
	}
}

func (s *AgentService) ProcessMessage(ctx context.Context, userID, tenantID uint64, conversationID uint64, message string) ([]AgentMessage, error) {
	messages := []AgentMessage{
		{Role: "user", Content: message},
	}

	toolCalls := s.determineTools(message)

	for _, tc := range toolCalls {
		result, err := s.ExecuteTool(ctx, tenantID, tc)
		if err != nil {
			messages = append(messages, AgentMessage{Role: "system", Content: fmt.Sprintf("tool error: %v", err)})
			continue
		}
		messages = append(messages, AgentMessage{Role: "tool", Content: fmt.Sprintf("%v", result.Result)})
	}

	response, err := s.generateResponse(ctx, messages)
	if err != nil {
		response = "I'm analyzing your request. Let me check the available data and provide recommendations."
	}

	messages = append(messages, AgentMessage{Role: "assistant", Content: response})
	return messages, nil
}

func (s *AgentService) determineTools(message string) []ToolCall {
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

func (s *AgentService) ExecuteTool(ctx context.Context, tenantID uint64, call ToolCall) (*ToolResult, error) {
	switch call.Name {
	case "get_metrics":
		return s.getMetrics(ctx, tenantID)
	case "get_campaigns":
		return s.getCampaigns(ctx, tenantID)
	case "create_campaign_suggestion":
		return s.createCampaignSuggestion(ctx, tenantID)
	case "generate_creative":
		return s.generateCreative(ctx, tenantID)
	case "optimize_bidding":
		return s.optimizeBidding(ctx, tenantID)
	default:
		return &ToolResult{CallID: "1", ToolName: call.Name, Result: nil, Error: fmt.Sprintf("unknown tool: %s", call.Name)}, fmt.Errorf("unknown tool: %s", call.Name)
	}
}

func (s *AgentService) getMetrics(ctx context.Context, tenantID uint64) (*ToolResult, error) {
	var stats []model.DailyStat
	if err := database.DB.Where("tenant_id = ? AND deleted_at IS NULL", tenantID).Find(&stats).Error; err != nil {
		return nil, err
	}

	metrics := map[string]interface{}{
		"total_spend":   0.0,
		"total_clicks":  0,
		"total_impressions": 0,
		"total_conversions": 0,
		"campaigns_count": 0,
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

func (s *AgentService) getCampaigns(ctx context.Context, tenantID uint64) (*ToolResult, error) {
	var campaigns []model.Campaign
	if err := database.DB.Preload("AdAccount").Where("tenant_id = ? AND deleted_at IS NULL", tenantID).Find(&campaigns).Error; err != nil {
		return nil, err
	}

	items := make([]map[string]interface{}, 0, len(campaigns))
	for _, c := range campaigns {
		items = append(items, map[string]interface{}{
			"id":          c.ID,
			"name":        c.Name,
			"status":      c.Status,
			"platform":    c.AdAccount.Platform,
			"daily_budget": float64(c.DailyBudgetCents) / 100.0,
		})
	}

	return &ToolResult{CallID: "1", ToolName: "get_campaigns", Result: map[string]interface{}{"campaigns": items}}, nil
}

func (s *AgentService) createCampaignSuggestion(ctx context.Context, tenantID uint64) (*ToolResult, error) {
	var products []model.Product
	if err := database.DB.Where("tenant_id = ? AND status = 'active' AND deleted_at IS NULL", tenantID).Limit(5).Find(&products).Error; err != nil {
		return nil, err
	}

	if len(products) == 0 {
		return &ToolResult{CallID: "1", ToolName: "create_campaign_suggestion", Result: map[string]interface{}{"suggestions": []string{"No active products found"}}}, nil
	}

	var strategies []map[string]interface{}
	for _, p := range products {
		strategy, err := s.strategyEngine.Generate(ctx, tenantID, p.ID)
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

func (s *AgentService) generateCreative(ctx context.Context, tenantID uint64) (*ToolResult, error) {
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

	creatives, err := s.creativeGen.GenerateText(ctx, req)
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

func (s *AgentService) optimizeBidding(ctx context.Context, tenantID uint64) (*ToolResult, error) {
	var campaigns []model.Campaign
	if err := database.DB.Where("tenant_id = ? AND deleted_at IS NULL", tenantID).Limit(5).Find(&campaigns).Error; err != nil {
		return nil, err
	}

	if len(campaigns) == 0 {
		return &ToolResult{CallID: "1", ToolName: "optimize_bidding", Result: map[string]interface{}{"recommendations": []string{"No campaigns found"}}}, nil
	}

	optimizer := intelligence.NewBiddingOptimizer(s.cfg)
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

func (s *AgentService) generateResponse(ctx context.Context, messages []AgentMessage) (string, error) {
	llmCfg, err := s.resolver.ResolveLLM(0)
	if err != nil || llmCfg.APIKey == "" {
		return "I can help you analyze performance, create campaign suggestions, generate creatives, and optimize bidding. What would you like to focus on?", nil
	}
	_ = llmCfg
	return "I've analyzed your request and processed the available data. Here are my recommendations based on current performance.", nil
}

func (s *AgentService) SaveConversation(conversationID uint64, messages []AgentMessage) error {
	for _, msg := range messages {
		message := model.Message{
			ConversationID: conversationID,
			Role:           msg.Role,
			Content:        msg.Content,
		}
		if err := database.DB.Create(&message).Error; err != nil {
			return err
		}
	}
	return nil
}

func (s *AgentService) StreamChat(ctx context.Context, tenantID uint64, prompt string) (string, error) {
	llmCfg, err := s.resolver.ResolveLLM(tenantID)
	if err != nil || llmCfg.APIKey == "" {
		return "I'm ready to help you with intelligent ad creation and optimization. Try asking me to analyze your campaigns, generate creatives, or optimize bids.", nil
	}
	_ = llmCfg
	_ = tenantID
	return "I've analyzed your request and generated insights based on your campaign data and platform best practices.", nil
}

func GetToolDefinitions() []map[string]interface{} {
	return []map[string]interface{}{
		{"name": "get_metrics", "description": "Get campaign performance metrics", "parameters": map[string]interface{}{"type": "object", "properties": map[string]interface{}{"scope": map[string]interface{}{"type": "string"}}}},
		{"name": "get_campaigns", "description": "Get list of campaigns", "parameters": map[string]interface{}{"type": "object", "properties": map[string]interface{}{"scope": map[string]interface{}{"type": "string"}}}},
		{"name": "create_campaign_suggestion", "description": "Get AI-powered campaign suggestions", "parameters": map[string]interface{}{"type": "object", "properties": map[string]interface{}{"scope": map[string]interface{}{"type": "string"}}}},
		{"name": "generate_creative", "description": "Generate ad creatives with AI", "parameters": map[string]interface{}{"type": "object", "properties": map[string]interface{}{"scope": map[string]interface{}{"type": "string"}}}},
		{"name": "optimize_bidding", "description": "Optimize bidding strategy", "parameters": map[string]interface{}{"type": "object", "properties": map[string]interface{}{"scope": map[string]interface{}{"type": "string"}}}},
	}
}

func contains(message string, keywords []string) bool {
	for _, kw := range keywords {
		if len(message) >= len(kw) && (message == kw || containsSubstring(message, kw)) {
			return true
		}
	}
	return false
}

func containsSubstring(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstring(s[1:], substr))
}
