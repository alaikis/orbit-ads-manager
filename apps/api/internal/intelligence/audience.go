package intelligence

import (
	"context"
	"fmt"
	"sort"
	"time"

	"orbit/apps/api/internal/model"
	"orbit/apps/api/pkg/database"
)

type AudienceAnalyzer struct{}

func NewAudienceAnalyzer() *AudienceAnalyzer {
	return &AudienceAnalyzer{}
}

func (a *AudienceAnalyzer) Analyze(ctx context.Context, tenantID uint64) ([]AudienceSegment, error) {
	var orders []model.Order
	if err := database.DB.Where("tenant_id = ? AND deleted_at IS NULL", tenantID).Find(&orders).Error; err != nil {
		return nil, fmt.Errorf("failed to load orders: %w", err)
	}

	var products []model.Product
	if err := database.DB.Where("tenant_id = ? AND deleted_at IS NULL", tenantID).Find(&products).Error; err != nil {
		return nil, fmt.Errorf("failed to load products: %w", err)
	}

	if len(orders) == 0 {
		return a.generateDefaultSegments(tenantID, products), nil
	}

	segments := []AudienceSegment{
		{
			ID:          fmt.Sprintf("seg_%d_high_value", tenantID),
			Name:        "High-Value Customers",
			Description: "Customers with total spend in top 20%",
			Size:       a.countHighValueCustomers(orders),
			Confidence:  0.85,
			Features: map[string]interface{}{
				"avg_order_value": a.avgOrderValue(orders),
				"frequency":       a.purchaseFrequency(orders),
				"recency":        a.recencyDays(orders),
			},
			Platforms: []string{"meta", "google"},
			CreatedAt: time.Now().UTC(),
		},
		{
			ID:          fmt.Sprintf("seg_%d_new", tenantID),
			Name:        "New Customers",
			Description: "First-time buyers in the last 30 days",
			Size:       a.countNewCustomers(orders),
			Confidence:  0.75,
			Features: map[string]interface{}{
				"avg_order_value": a.avgOrderValue(orders),
				"frequency":       1.0,
				"recency":         a.recencyDays(orders),
			},
			Platforms: []string{"meta", "google"},
			CreatedAt: time.Now().UTC(),
		},
		{
			ID:          fmt.Sprintf("seg_%d_product_affinity", tenantID),
			Name:        "Product Affinity Groups",
			Description: "Customers grouped by product preferences",
			Size:       len(orders),
			Confidence:  0.70,
			Features: map[string]interface{}{
				"top_categories": a.topCategories(orders, products),
			},
			Platforms: []string{"meta", "google"},
			CreatedAt: time.Now().UTC(),
		},
	}

	sort.Slice(segments, func(i, j int) bool {
		return segments[i].Size > segments[j].Size
	})

	return segments, nil
}

func (a *AudienceAnalyzer) generateDefaultSegments(tenantID uint64, products []model.Product) []AudienceSegment {
	categories := make([]string, 0, len(products))
	for _, p := range products {
		if p.ProductType != "" {
			categories = append(categories, p.ProductType)
		}
	}
	if len(categories) == 0 {
		categories = append(categories, "general")
	}

	return []AudienceSegment{
		{
			ID:          fmt.Sprintf("seg_%d_general", tenantID),
			Name:        "General Audience",
			Description: "Broad targeting based on store category",
			Size:       1000,
			Confidence:  0.50,
			Features: map[string]interface{}{
				"categories": categories,
			},
			Platforms: []string{"meta", "google"},
			CreatedAt: time.Now().UTC(),
		},
	}
}

func (a *AudienceAnalyzer) countHighValueCustomers(orders []model.Order) int {
	if len(orders) == 0 {
		return 0
	}

	customerTotals := make(map[string]int64)
	for _, o := range orders {
		email := o.CustomerEmail
		if email == "" {
			email = fmt.Sprintf("order_%d", o.ID)
		}
		customerTotals[email] += o.TotalCents
	}

	totals := make([]int64, 0, len(customerTotals))
	for _, total := range customerTotals {
		totals = append(totals, total)
	}

	if len(totals) == 0 {
		return 0
	}

	sort.Slice(totals, func(i, j int) bool { return totals[i] > totals[j] })
	thresholdIdx := int(float64(len(totals)) * 0.8)
	if thresholdIdx >= len(totals) {
		thresholdIdx = len(totals) - 1
	}

	count := 0
	threshold := totals[thresholdIdx]
	for _, total := range customerTotals {
		if total >= threshold {
			count++
		}
	}

	return count
}

func (a *AudienceAnalyzer) countNewCustomers(orders []model.Order) int {
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	count := 0
	customerSeen := make(map[string]bool)

	for _, o := range orders {
		if o.PlacedAt.After(thirtyDaysAgo) && !customerSeen[o.CustomerEmail] {
			customerSeen[o.CustomerEmail] = true
			count++
		}
	}

	return count
}

func (a *AudienceAnalyzer) avgOrderValue(orders []model.Order) float64 {
	if len(orders) == 0 {
		return 0
	}
	total := int64(0)
	for _, o := range orders {
		total += o.TotalCents
	}
	return float64(total) / float64(len(orders)) / 100.0
}

func (a *AudienceAnalyzer) purchaseFrequency(orders []model.Order) float64 {
	if len(orders) == 0 {
		return 0
	}
	customerCount := make(map[string]bool)
	for _, o := range orders {
		if o.CustomerEmail != "" {
			customerCount[o.CustomerEmail] = true
		}
	}
	return float64(len(orders)) / float64(len(customerCount))
}

func (a *AudienceAnalyzer) recencyDays(orders []model.Order) float64 {
	if len(orders) == 0 {
		return 0
	}
	mostRecent := orders[0].PlacedAt
	for _, o := range orders {
		if o.PlacedAt.After(mostRecent) {
			mostRecent = o.PlacedAt
		}
	}
	return time.Since(mostRecent).Hours() / 24.0
}

func (a *AudienceAnalyzer) topCategories(orders []model.Order, products []model.Product) []string {
	categoryCount := make(map[string]int)
	productMap := make(map[uint64]string)
	for _, p := range products {
		productMap[p.ID] = p.ProductType
	}

	for _, o := range orders {
		for _, item := range o.Items {
			productIDFloat, ok := item["product_id"].(float64)
			if !ok {
				continue
			}
			productID := uint64(productIDFloat)
			if cat, ok := productMap[productID]; ok && cat != "" {
				categoryCount[cat]++
			}
		}
	}

	type catCount struct {
		category string
		count    int
	}
	cats := make([]catCount, 0, len(categoryCount))
	for cat, count := range categoryCount {
		cats = append(cats, catCount{cat, count})
	}
	sort.Slice(cats, func(i, j int) bool { return cats[i].count > cats[j].count })

	result := make([]string, 0, len(cats))
	for _, c := range cats {
		result = append(result, c.category)
	}
	return result
}

func (a *AudienceAnalyzer) AnalyzeProductMargins(ctx context.Context, tenantID uint64, products []model.Product) ([]ProductInsight, error) {
	if len(products) == 0 {
		return nil, nil
	}

	var orders []model.Order
	if err := database.DB.Where("tenant_id = ? AND deleted_at IS NULL", tenantID).Find(&orders).Error; err != nil {
		return nil, err
	}

	productSales := make(map[uint64]int)
	for _, o := range orders {
		for _, item := range o.Items {
			productIDFloat, ok := item["product_id"].(float64)
			if !ok {
				continue
			}
			productID := uint64(productIDFloat)
			qtyFloat, ok := item["quantity"].(float64)
			if !ok {
				continue
			}
			productSales[productID] += int(qtyFloat)
		}
	}

	insights := make([]ProductInsight, 0, len(products))
	for _, p := range products {
		if p.Status != "active" {
			continue
		}

		popularity := productSales[p.ID]
		margin := 0.35
		if p.ComparePriceCents != nil && *p.ComparePriceCents > 0 {
			margin = float64(p.PriceCents-*p.ComparePriceCents) / float64(p.PriceCents)
		}
		if margin < 0 {
			margin = 0.10
		}
		if margin > 0.8 {
			margin = 0.8
		}

		bestPlatform := "meta"
		if p.ProductType == "electronics" || p.ProductType == "gadgets" {
			bestPlatform = "google"
		} else if p.ProductType == "fashion" || p.ProductType == "apparel" {
			bestPlatform = "meta"
		}

		bid := float64(p.PriceCents) / 100.0 * margin * 0.3
		if bid < 0.5 {
			bid = 0.5
		}

		insights = append(insights, ProductInsight{
			ProductID:      p.ID,
			ProductName:    p.Title,
			Margin:         mathRound(margin*100) / 100,
			Popularity:     popularity,
			Seasonality:    detectSeasonality(p),
			BestPlatform:   bestPlatform,
			BestAudience:   fmt.Sprintf("seg_%d_product_affinity", tenantID),
			RecommendedBid: mathRound(bid*100) / 100,
		})
	}

	sort.Slice(insights, func(i, j int) bool {
		return insights[i].Popularity > insights[j].Popularity
	})

	return insights, nil
}

func detectSeasonality(product model.Product) string {
	if product.GoogleCategory != "" {
		switch product.GoogleCategory {
		case "Seasonal", "Holiday", "Christmas", "Halloween":
			return "seasonal_peak"
		}
	}
	return "year_round"
}

func mathRound(x float64) float64 {
	return float64(int64(x*100+0.5)) / 100
}
