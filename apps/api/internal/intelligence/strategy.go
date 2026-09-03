package intelligence

import (
	"context"
	"fmt"
	"math"
	"sort"

	"orbit/apps/api/internal/model"
	"orbit/apps/api/pkg/database"
	"orbit/apps/api/config"
)

type StrategyEngine struct {
	analyzer *AudienceAnalyzer
	cfg      *config.Config
}

func NewStrategyEngine(cfg *config.Config) *StrategyEngine {
	return &StrategyEngine{
		analyzer: NewAudienceAnalyzer(),
		cfg:      cfg,
	}
}

type PlatformProfile struct {
	Platform          string  `json:"platform"`
	Strengths         []string `json:"strengths"`
	BestFor           []string `json:"best_for"`
	AvgCPC            float64 `json:"avg_cpc"`
	AvgCTR            float64 `json:"avg_ctr"`
	AvgConversionRate float64 `json:"avg_conversion_rate"`
}

var PlatformProfiles = map[string]PlatformProfile{
	"meta": {
		Platform:          "meta",
		Strengths:         []string{"visual_products", "brand_awareness", "retargeting", "social_proof"},
		BestFor:           []string{"fashion", "beauty", "food", "travel", "entertainment"},
		AvgCPC:            0.50,
		AvgCTR:            0.012,
		AvgConversionRate: 0.025,
	},
	"google": {
		Platform:          "google",
		Strengths:         []string{"high_intent", "search_ads", "shopping", "local"},
		BestFor:           []string{"electronics", "home_garden", "automotive", "health"},
		AvgCPC:            0.80,
		AvgCTR:            0.035,
		AvgConversionRate: 0.035,
	},
}

func (e *StrategyEngine) Generate(ctx context.Context, tenantID uint64, productID uint64) (*CampaignStrategy, error) {
	var product model.Product
	if err := database.DB.Where("id = ? AND tenant_id = ?", productID, tenantID).First(&product).Error; err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	segments, err := e.analyzer.Analyze(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	insights, err := e.analyzer.AnalyzeProductMargins(ctx, tenantID, []model.Product{product})
	if err != nil {
		return nil, err
	}

	if len(insights) == 0 {
		return nil, fmt.Errorf("no product insights available")
	}

	insight := insights[0]
	profile := PlatformProfiles[insight.BestPlatform]

	strategy := &CampaignStrategy{
		Platform:    insight.BestPlatform,
		Objective:   "conversions",
		BudgetSplit: 0.6,
		Audience: []string{
			fmt.Sprintf("High-Value Customers (segment)"),
			fmt.Sprintf("Product Affinity - %s", product.ProductType),
		},
		BiddingType: "bid_target",
		BidAmount:   insight.RecommendedBid,
		Creatives: []CreativeSuggestion{
			{
				Type:        "image",
				Headline:    generateHeadline(product.Title, product.PriceCents, product.Currency),
				Description: generateDescription(product),
				ImagePrompt: generateImagePrompt(product),
				CTA:         "Shop Now",
			},
		},
		Schedule: []TimeSlot{
			{DayOfWeek: 0, HourStart: 9, HourEnd: 21, BidAdjust: 1.1},
			{DayOfWeek: 1, HourStart: 8, HourEnd: 22, BidAdjust: 1.2},
			{DayOfWeek: 5, HourStart: 10, HourEnd: 23, BidAdjust: 1.3},
		},
		Reasoning: fmt.Sprintf(
			"Platform %s selected based on product category (%s) and audience affinity. "+
			"Bid optimized for %.2fx margin with target ROAS. "+
			"Creative emphasizes %s strengths: %s.",
			insight.BestPlatform,
			product.ProductType,
			insight.Margin,
			insight.BestPlatform,
			fmt.Sprintf("%v", profile.Strengths[:min(2, len(profile.Strengths))]),
		),
	}

	for _, seg := range segments {
		if seg.Name == "High-Value Customers" {
			strategy.Audience = append(strategy.Audience, seg.ID)
		}
	}

	return strategy, nil
}

func (e *StrategyEngine) GenerateMultiPlatform(ctx context.Context, tenantID uint64) ([]CampaignStrategy, error) {
	var products []model.Product
	if err := database.DB.Where("tenant_id = ? AND status = 'active' AND deleted_at IS NULL", tenantID).Find(&products).Error; err != nil {
		return nil, err
	}

	if len(products) == 0 {
		return nil, fmt.Errorf("no active products")
	}

	insights, err := e.analyzer.AnalyzeProductMargins(ctx, tenantID, products)
	if err != nil {
		return nil, err
	}

	platformProducts := make(map[string][]ProductInsight)
	for _, ins := range insights {
		platformProducts[ins.BestPlatform] = append(platformProducts[ins.BestPlatform], ins)
	}

	var strategies []CampaignStrategy
	for platform, prods := range platformProducts {
		profile := PlatformProfiles[platform]
		totalMargin := 0.0
		for _, p := range prods {
			totalMargin += p.Margin
		}
		avgMargin := totalMargin / float64(len(prods))

		strategy := CampaignStrategy{
			Platform:    platform,
			Objective:   "conversions",
			BudgetSplit: 1.0 / float64(len(platformProducts)),
			Audience: []string{
				fmt.Sprintf("High-Value Customers"),
				fmt.Sprintf("Product Affinity"),
			},
			BiddingType: "bid_target",
			BidAmount:   math.Round(avgMargin*50*100) / 100,
			Creatives:   make([]CreativeSuggestion, 0, len(prods)),
			Schedule: []TimeSlot{
				{DayOfWeek: 0, HourStart: 9, HourEnd: 21, BidAdjust: 1.0},
				{DayOfWeek: 1, HourStart: 8, HourEnd: 22, BidAdjust: 1.15},
				{DayOfWeek: 5, HourStart: 10, HourEnd: 23, BidAdjust: 1.2},
			},
			Reasoning: fmt.Sprintf(
				"Multi-platform strategy for %d products on %s. "+
				"Avg margin: %.2f. Platform strengths: %s.",
				len(prods), platform, avgMargin, fmt.Sprintf("%v", profile.Strengths[:min(2, len(profile.Strengths))]),
			),
		}

		for _, prod := range prods {
			strategy.Creatives = append(strategy.Creatives, CreativeSuggestion{
				Type:        "image",
				Headline:    generateHeadline(prod.ProductName, int64(prod.RecommendedBid*100), "USD"),
				Description: fmt.Sprintf("Popular product with %.0f%% margin. Shop now!", prod.Margin*100),
				ImagePrompt: fmt.Sprintf("Professional product photo of %s, clean background, e-commerce style", prod.ProductName),
				CTA:         "Shop Now",
			})
		}

		strategies = append(strategies, strategy)
	}

	sort.Slice(strategies, func(i, j int) bool {
		return strategies[i].BudgetSplit > strategies[j].BudgetSplit
	})

	return strategies, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func generateHeadline(title string, priceCents int64, currency string) string {
	price := float64(priceCents) / 100.0
	return fmt.Sprintf("%s - $%.2f", title, price)
}

func generateDescription(product model.Product) string {
	if product.Description != "" {
		if len(product.Description) > 120 {
			return product.Description[:117] + "..."
		}
		return product.Description
	}
	return fmt.Sprintf("Shop %s today! Free shipping on orders over $50.", product.Title)
}

func generateImagePrompt(product model.Product) string {
	if product.Brand != "" {
		return fmt.Sprintf("Professional product photography of %s by %s, clean white background, high resolution, e-commerce style", product.Title, product.Brand)
	}
	return fmt.Sprintf("Professional product photography of %s, clean white background, high resolution, e-commerce style", product.Title)
}

func (e *StrategyEngine) OptimizeBudget(ctx context.Context, tenantID uint64, totalBudget float64) (map[string]float64, error) {
	var dailyStats []model.DailyStat
	if err := database.DB.Where("tenant_id = ? AND deleted_at IS NULL", tenantID).Find(&dailyStats).Error; err != nil {
		return nil, fmt.Errorf("failed to load stats: %w", err)
	}

	platformPerformance := make(map[string]float64)
	for _, stat := range dailyStats {
		spend, _ := stat.Metrics["spend"].(float64)
		revenue, _ := stat.Metrics["revenue"].(float64)
		if spend > 0 {
			roas := revenue / spend
			platformPerformance[stat.Platform] += roas
		}
	}

	totalPerformance := 0.0
	for _, perf := range platformPerformance {
		totalPerformance += perf
	}

	if totalPerformance == 0 {
		return map[string]float64{
			"meta":   totalBudget * 0.5,
			"google": totalBudget * 0.5,
		}, nil
	}

	budgetSplit := make(map[string]float64)
	for platform, perf := range platformPerformance {
		share := perf / totalPerformance
		budgetSplit[platform] = math.Round(totalBudget*share*100) / 100
	}

	remaining := totalBudget
	for _, split := range budgetSplit {
		remaining -= split
	}

	if remaining > 0 {
		for platform := range budgetSplit {
			budgetSplit[platform] += remaining / float64(len(budgetSplit))
			break
		}
	}

	return budgetSplit, nil
}
