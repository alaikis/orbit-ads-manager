package intelligence

import (
	"context"
	"fmt"
	"math"
	"time"

	"orbit/apps/api/internal/model"
	"orbit/apps/api/pkg/database"
	"orbit/apps/api/config"
)

type BiddingOptimizer struct {
	cfg *config.Config
}

func NewBiddingOptimizer(cfg *config.Config) *BiddingOptimizer {
	return &BiddingOptimizer{cfg: cfg}
}

type BiddingSignal struct {
	Metric       string  `json:"metric"`
	Value        float64 `json:"value"`
	Weight       float64 `json:"weight"`
	Direction    string  `json:"direction"`
}

func (o *BiddingOptimizer) Optimize(ctx context.Context, tenantID uint64, campaignID uint64) (*BiddingRecommendation, error) {
	var campaign model.Campaign
	if err := database.DB.Preload("AdAccount").Where("id = ? AND tenant_id = ?", campaignID, tenantID).First(&campaign).Error; err != nil {
		return nil, fmt.Errorf("campaign not found: %w", err)
	}

	var dailyStats []model.DailyStat
	if err := database.DB.Where("tenant_id = ? AND scope_type = 'campaign' AND scope_id = ? AND date >= ? AND deleted_at IS NULL", tenantID, campaign.ID, time.Now().AddDate(0, -1, 0).Format("2006-01-02")).Find(&dailyStats).Error; err != nil {
		return nil, fmt.Errorf("failed to load daily stats: %w", err)
	}

	var account model.AdAccount
	if err := database.DB.Where("id = ?", campaign.AdAccountID).First(&account).Error; err != nil {
		return nil, fmt.Errorf("ad account not found: %w", err)
	}

	var conn model.PlatformConn
	if err := database.DB.Where("id = ?", account.ConnID).First(&conn).Error; err != nil {
		return nil, fmt.Errorf("platform connection not found: %w", err)
	}

	var oauthToken model.OAuthToken
	if err := database.DB.Where("conn_id = ?", conn.ID).First(&oauthToken).Error; err != nil {
		return nil, fmt.Errorf("oauth token not found: %w", err)
	}

	_, currentBid := o.calculateCurrentBid(campaign)

	signals := o.collectSignals(dailyStats, campaign)
	recommended := o.calculateOptimalBid(signals, currentBid)

	reason := o.generateReason(signals, recommended, currentBid)

	return &BiddingRecommendation{
		CampaignID:  campaign.ID,
		CurrentBid:  currentBid,
		Recommended: math.Round(recommended*100) / 100,
		Min:         math.Round(currentBid*0.5*100) / 100,
		Max:         math.Round(currentBid*2.0*100) / 100,
		Reason:      reason,
		Confidence:  o.calculateConfidence(signals),
		UpdatedAt:   time.Now().UTC(),
	}, nil
}

func (o *BiddingOptimizer) collectSignals(stats []model.DailyStat, campaign model.Campaign) []BiddingSignal {
	signals := []BiddingSignal{}

	if len(stats) == 0 {
		return signals
	}

	totalSpend := 0.0
	totalClicks := 0.0
	totalImpressions := 0.0
	totalConversions := 0.0

	for _, stat := range stats {
		spend, _ := stat.Metrics["spend"].(float64)
		clicks, _ := stat.Metrics["clicks"].(float64)
		impressions, _ := stat.Metrics["impressions"].(float64)
		conversions, _ := stat.Metrics["conversions"].(float64)

		totalSpend += spend
		totalClicks += clicks
		totalImpressions += impressions
		totalConversions += conversions
	}

	if totalImpressions > 0 {
		ctr := totalClicks / totalImpressions
		signals = append(signals, BiddingSignal{
			Metric:    "ctr",
			Value:     ctr,
			Weight:    0.2,
			Direction: "higher_is_better",
		})
	}

	if totalClicks > 0 {
		cpc := totalSpend / totalClicks
		signals = append(signals, BiddingSignal{
			Metric:    "cpc",
			Value:     cpc,
			Weight:    0.3,
			Direction: "lower_is_better",
		})
	}

	if totalConversions > 0 && totalSpend > 0 {
		cpa := totalSpend / totalConversions
		signals = append(signals, BiddingSignal{
			Metric:    "cpa",
			Value:     cpa,
			Weight:    0.4,
			Direction: "lower_is_better",
		})
	}

	if totalConversions > 0 && totalClicks > 0 {
		cvr := totalConversions / totalClicks
		signals = append(signals, BiddingSignal{
			Metric:    "conversion_rate",
			Value:     cvr,
			Weight:    0.3,
			Direction: "higher_is_better",
		})
	}

	if campaign.DailyBudgetCents > 0 {
		avgDailySpend := totalSpend / float64(len(stats))
		budgetUtilization := avgDailySpend / (float64(campaign.DailyBudgetCents) / 100.0)
		signals = append(signals, BiddingSignal{
			Metric:    "budget_utilization",
			Value:     budgetUtilization,
			Weight:    0.2,
			Direction: "optimal_0.8",
		})
	}

	return signals
}

func (o *BiddingOptimizer) calculateCurrentBid(campaign model.Campaign) (float64, float64) {
	if campaign.DailyBudgetCents > 0 {
		return float64(campaign.DailyBudgetCents) / 100.0, float64(campaign.DailyBudgetCents) / 100.0
	}
	return 1.0, 1.0
}

func (o *BiddingOptimizer) calculateOptimalBid(signals []BiddingSignal, currentBid float64) float64 {
	if len(signals) == 0 {
		return currentBid
	}

	adjustment := 0.0
	totalWeight := 0.0

	for _, signal := range signals {
		totalWeight += signal.Weight

		switch signal.Direction {
		case "higher_is_better":
			if signal.Value > 0.02 {
				adjustment += signal.Weight * 0.1
			} else if signal.Value < 0.005 {
				adjustment -= signal.Weight * 0.15
			}
		case "lower_is_better":
			if signal.Value < 0.5 {
				adjustment += signal.Weight * 0.1
			} else if signal.Value > 1.5 {
				adjustment -= signal.Weight * 0.15
			}
		case "optimal_0.8":
			if signal.Value < 0.5 {
				adjustment += signal.Weight * 0.15
			} else if signal.Value > 0.95 {
				adjustment -= signal.Weight * 0.1
			}
		}
	}

	if totalWeight > 0 {
		adjustment /= totalWeight
	}

	adjustment = math.Max(-0.3, math.Min(0.3, adjustment))

	optimal := currentBid * (1 + adjustment)
	return math.Max(0.1, math.Min(currentBid*2.0, optimal))
}

func (o *BiddingOptimizer) calculateConfidence(signals []BiddingSignal) float64 {
	if len(signals) == 0 {
		return 0.3
	}

	dataPoints := len(signals)
	confidence := math.Min(0.95, 0.3+float64(dataPoints)*0.1)

	return math.Round(confidence*100) / 100
}

func (o *BiddingOptimizer) generateReason(signals []BiddingSignal, recommended, currentBid float64) string {
	if len(signals) == 0 {
		return "Insufficient data for optimization"
	}

	if recommended > currentBid*1.1 {
		return "Strong performance signals detected. Increasing bid to capture more traffic."
	} else if recommended < currentBid*0.9 {
		return "Weak performance signals detected. Decreasing bid to improve efficiency."
	}
	return "Performance is stable. Maintaining current bid with minor adjustments."
}
