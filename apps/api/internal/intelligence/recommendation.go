package intelligence

import (
	"context"
	"fmt"
	"math"
	"sort"
	"time"

	"orbit/apps/api/internal/model"
	"orbit/apps/api/pkg/database"
	"orbit/apps/api/config"
)

type RecommendationEngine struct {
	analyzer      *AudienceAnalyzer
	strategyEngine *StrategyEngine
	biddingOptimizer *BiddingOptimizer
}

func NewRecommendationEngine(cfg *config.Config) *RecommendationEngine {
	return &RecommendationEngine{
		analyzer:       NewAudienceAnalyzer(),
		strategyEngine: NewStrategyEngine(cfg),
		biddingOptimizer: NewBiddingOptimizer(cfg),
	}
}

type CampaignRecommendation struct {
	ID            string    `json:"id"`
	ProductID     uint64    `json:"product_id"`
	ProductName   string    `json:"product_name"`
	Platform      string    `json:"platform"`
	Strategy      CampaignStrategy `json:"strategy"`
	Bid           BiddingRecommendation `json:"bid"`
	ExpectedROAS  float64   `json:"expected_roas"`
	Confidence    float64   `json:"confidence"`
	Priority      int       `json:"priority"`
	Reasoning     string    `json:"reasoning"`
	CreatedAt     time.Time `json:"created_at"`
}

func (e *RecommendationEngine) GenerateRecommendations(ctx context.Context, tenantID uint64) ([]CampaignRecommendation, error) {
	var products []model.Product
	if err := database.DB.Where("tenant_id = ? AND status = 'active' AND deleted_at IS NULL", tenantID).Find(&products).Error; err != nil {
		return nil, fmt.Errorf("failed to load products: %w", err)
	}

	if len(products) == 0 {
		return nil, fmt.Errorf("no active products found")
	}

	insights, err := e.analyzer.AnalyzeProductMargins(ctx, tenantID, products)
	if err != nil {
		return nil, err
	}

	var recommendations []CampaignRecommendation
	for _, insight := range insights {
		strategy, err := e.strategyEngine.Generate(ctx, tenantID, insight.ProductID)
		if err != nil {
			continue
		}

		campaignID := uint64(insight.ProductID)
		bid, err := e.biddingOptimizer.Optimize(ctx, tenantID, campaignID)
		if err != nil {
			bid = &BiddingRecommendation{
				CampaignID:  campaignID,
				CurrentBid:  insight.RecommendedBid,
				Recommended: insight.RecommendedBid,
				Min:         insight.RecommendedBid * 0.5,
				Max:         insight.RecommendedBid * 2.0,
				Reason:      "Insufficient data for optimization",
				Confidence:  0.4,
				UpdatedAt:   time.Now().UTC(),
			}
		}

		expectedROAS := e.estimateROAS(insight, strategy, bid)

		priority := e.calculatePriority(insight, expectedROAS, bid.Confidence)

		rec := CampaignRecommendation{
			ID:          fmt.Sprintf("rec_%d_%d", tenantID, insight.ProductID),
			ProductID:   insight.ProductID,
			ProductName: insight.ProductName,
			Platform:    strategy.Platform,
			Strategy:    *strategy,
			Bid:         *bid,
			ExpectedROAS: expectedROAS,
			Confidence:  bid.Confidence,
			Priority:    priority,
			Reasoning:   strategy.Reasoning,
			CreatedAt:   time.Now().UTC(),
		}

		recommendations = append(recommendations, rec)
	}

	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].Priority > recommendations[j].Priority
	})

	return recommendations, nil
}

func (e *RecommendationEngine) estimateROAS(insight ProductInsight, strategy *CampaignStrategy, bid *BiddingRecommendation) float64 {
	platformProfile := PlatformProfiles[strategy.Platform]
	baseROAS := platformProfile.AvgConversionRate * platformProfile.AvgCTR * 100.0
	marginBoost := insight.Margin * 2.0
	roas := baseROAS * (1 + marginBoost)

	if bid.Recommended < bid.CurrentBid*0.8 {
		roas *= 1.1
	} else if bid.Recommended > bid.CurrentBid*1.2 {
		roas *= 0.9
	}

	return math.Round(roas*100) / 100
}

func (e *RecommendationEngine) calculatePriority(insight ProductInsight, expectedROAS float64, confidence float64) int {
	score := 0

	if expectedROAS > 3.0 {
		score += 40
	} else if expectedROAS > 2.0 {
		score += 30
	} else if expectedROAS > 1.0 {
		score += 20
	}

	if insight.Popularity > 50 {
		score += 20
	} else if insight.Popularity > 10 {
		score += 10
	}

	if confidence > 0.7 {
		score += 20
	} else if confidence > 0.5 {
		score += 10
	}

	if insight.Margin > 0.4 {
		score += 20
	} else if insight.Margin > 0.2 {
		score += 10
	}

	return score
}

func (e *RecommendationEngine) AnalyzeCompetition(ctx context.Context, tenantID uint64) (map[string]interface{}, error) {
	var products []model.Product
	if err := database.DB.Where("tenant_id = ? AND status = 'active' AND deleted_at IS NULL", tenantID).Find(&products).Error; err != nil {
		return nil, fmt.Errorf("failed to load products: %w", err)
	}

	analysis := map[string]interface{}{
		"total_products": len(products),
		"platforms":     make(map[string]int),
		"avg_price":     0.0,
		"avg_margin":    0.0,
		"competition_level": "medium",
	}

	if len(products) == 0 {
		return analysis, nil
	}

	totalPrice := int64(0)
	for _, p := range products {
		totalPrice += p.PriceCents
		if p.ProductType != "" {
			analysis["platforms"].(map[string]int)[PlatformProfiles["meta"].BestFor[0]]++
		}
	}
	analysis["avg_price"] = math.Round(float64(totalPrice)/float64(len(products))/100.0*100) / 100

	return analysis, nil
}
