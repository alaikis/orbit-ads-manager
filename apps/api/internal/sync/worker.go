package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"orbit/apps/api/internal/model"
	"orbit/apps/api/internal/feed"
	"orbit/apps/api/internal/rule"
	"orbit/apps/api/internal/store"
	"orbit/apps/api/pkg/database"
	"orbit/apps/api/config"

	"github.com/hibiken/asynq"
)

func RegisterHandlers(mux *asynq.ServeMux) {
	mux.HandleFunc("sync:products", handleSyncProducts)
	mux.HandleFunc("sync:orders", handleSyncOrders)
	mux.HandleFunc("sync:ad_meta", handleSyncAdMeta)
	mux.HandleFunc("sync:ad_metrics", handleSyncAdMetrics)
	mux.HandleFunc("feed:regenerate", handleFeedRegenerate)
	mux.HandleFunc("rule:evaluate", handleRuleEvaluate)
}

func handleSyncProducts(ctx context.Context, task *asynq.Task) error {
	var payload struct {
		TenantID uint64 `json:"tenant_id"`
		StoreID  uint64 `json:"store_id"`
		Type     string `json:"type"`
		Mode     string `json:"mode"`
	}
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return err
	}

	job := model.SyncJob{
		BaseModel: model.BaseModel{TenantID: payload.TenantID},
		Type:      payload.Type,
		Scope:     "store",
		ObjectID:  &payload.StoreID,
		Status:    "running",
		StartedAt: func() *time.Time { t := time.Now().UTC(); return &t }(),
	}
	database.DB.Create(&job)

	var conn model.PlatformConn
	if err := database.DB.Where("id = ? AND tenant_id = ?", payload.StoreID, payload.TenantID).First(&conn).Error; err != nil {
		job.Status = "failed"
		job.ErrorMsg = fmt.Sprintf("store not found: %v", err)
		database.DB.Save(&job)
		return fmt.Errorf("store not found: %w", err)
	}

	var storeConfig model.StoreConfig
	if err := database.DB.Where("store_id = ?", conn.ID).First(&storeConfig).Error; err != nil {
		job.Status = "failed"
		job.ErrorMsg = fmt.Sprintf("store config not found: %v", err)
		database.DB.Save(&job)
		return fmt.Errorf("store config not found: %w", err)
	}

	apiKey, _ := storeConfig.ExtraHeaders["api_key"].(string)
	apiSecret, _ := storeConfig.ExtraHeaders["api_secret"].(string)

	cfg := config.AppCfg
	var products []store.Product
	var err error

	switch conn.Platform {
	case "woocommerce":
		adapter := store.NewWooCommerceAdapter(cfg)
		products, err = adapter.SyncProductsFromAPI(conn.ID, storeConfig.BaseURL, apiKey, apiSecret)
	case "shopify":
		adapter := store.NewShopifyAdapter(cfg)
		products, err = adapter.SyncProductsFromAPI(conn.ID, storeConfig.BaseURL, apiKey)
	default:
		job.Status = "failed"
		job.ErrorMsg = fmt.Sprintf("unsupported platform: %s", conn.Platform)
		database.DB.Save(&job)
		return fmt.Errorf("unsupported platform: %s", conn.Platform)
	}

	if err != nil {
		job.Status = "failed"
		job.ErrorMsg = err.Error()
		database.DB.Save(&job)
		return err
	}

	for _, p := range products {
		product := model.Product{
			BaseModel:  model.BaseModel{TenantID: payload.TenantID},
			StoreID:    payload.StoreID,
			ExternalID: p.ID,
			Title:      p.Title,
			Description: p.Description,
			Link:       p.Link,
			ImageURL:   p.ImageURL,
			Brand:      p.Brand,
			Gtin:       p.Gtin,
			PriceCents: p.PriceCents,
			Status:     p.Status,
		}

		var existing model.Product
		if err := database.DB.Where("tenant_id = ? AND store_id = ? AND external_id = ?", payload.TenantID, payload.StoreID, p.ID).First(&existing).Error; err == nil {
			existing.Title = p.Title
			existing.Description = p.Description
			existing.Link = p.Link
			existing.ImageURL = p.ImageURL
			existing.Brand = p.Brand
			existing.Gtin = p.Gtin
			existing.PriceCents = p.PriceCents
			existing.Status = p.Status
			database.DB.Save(&existing)
		} else {
			database.DB.Create(&product)
		}
	}

	now := time.Now().UTC()
	job.Status = "success"
	job.FinishedAt = &now
	job.RetryCount = 0
	database.DB.Save(&job)

	conn.LastSyncedAt = &now
	database.DB.Save(&conn)

	return nil
}

func handleSyncOrders(ctx context.Context, task *asynq.Task) error {
	var payload struct {
		TenantID uint64 `json:"tenant_id"`
		StoreID  uint64 `json:"store_id"`
		Type     string `json:"type"`
		Since    string `json:"since"`
	}
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return err
	}

	job := model.SyncJob{
		BaseModel: model.BaseModel{TenantID: payload.TenantID},
		Type:      payload.Type,
		Scope:     "store",
		ObjectID:  &payload.StoreID,
		Status:    "running",
		StartedAt: func() *time.Time { t := time.Now().UTC(); return &t }(),
	}
	database.DB.Create(&job)

	var conn model.PlatformConn
	if err := database.DB.Where("id = ? AND tenant_id = ?", payload.StoreID, payload.TenantID).First(&conn).Error; err != nil {
		job.Status = "failed"
		job.ErrorMsg = fmt.Sprintf("store not found: %v", err)
		database.DB.Save(&job)
		return fmt.Errorf("store not found: %w", err)
	}

	var storeConfig model.StoreConfig
	if err := database.DB.Where("store_id = ?", conn.ID).First(&storeConfig).Error; err != nil {
		job.Status = "failed"
		job.ErrorMsg = fmt.Sprintf("store config not found: %v", err)
		database.DB.Save(&job)
		return fmt.Errorf("store config not found: %w", err)
	}

	apiKey, _ := storeConfig.ExtraHeaders["api_key"].(string)
	apiSecret, _ := storeConfig.ExtraHeaders["api_secret"].(string)

	cfg := config.AppCfg
	var orders []store.Order
	var err error

	switch conn.Platform {
	case "woocommerce":
		adapter := store.NewWooCommerceAdapter(cfg)
		orders, err = adapter.SyncOrdersFromAPI(conn.ID, storeConfig.BaseURL, apiKey, apiSecret, payload.Since)
	case "shopify":
		adapter := store.NewShopifyAdapter(cfg)
		orders, err = adapter.SyncOrdersFromAPI(conn.ID, storeConfig.BaseURL, apiKey, payload.Since)
	default:
		job.Status = "failed"
		job.ErrorMsg = fmt.Sprintf("unsupported platform: %s", conn.Platform)
		database.DB.Save(&job)
		return fmt.Errorf("unsupported platform: %s", conn.Platform)
	}

	if err != nil {
		job.Status = "failed"
		job.ErrorMsg = err.Error()
		database.DB.Save(&job)
		return err
	}

	for _, o := range orders {
		placedAt := time.Now().UTC()
		if o.PlacedAt != "" {
			if parsed, parseErr := time.Parse(time.RFC3339, o.PlacedAt); parseErr == nil {
				placedAt = parsed
			}
		}

		order := model.Order{
			BaseModel:     model.BaseModel{TenantID: payload.TenantID},
			StoreID:       payload.StoreID,
			ExternalID:    o.ID,
			OrderNumber:   o.OrderNumber,
			Status:        o.Status,
			TotalCents:    o.TotalCents,
			Currency:      o.Currency,
			CustomerEmail: o.CustomerEmail,
			Items:         o.Items,
			PlacedAt:      placedAt,
		}

		var existing model.Order
		if err := database.DB.Where("tenant_id = ? AND store_id = ? AND external_id = ?", payload.TenantID, payload.StoreID, o.ID).First(&existing).Error; err == nil {
			existing.Status = o.Status
			existing.TotalCents = o.TotalCents
			existing.Items = o.Items
			database.DB.Save(&existing)
		} else {
			database.DB.Create(&order)
		}
	}

	now := time.Now().UTC()
	job.Status = "success"
	job.FinishedAt = &now
	job.RetryCount = 0
	database.DB.Save(&job)

	return nil
}

func handleSyncAdMeta(ctx context.Context, task *asynq.Task) error {
	return nil
}

func handleSyncAdMetrics(ctx context.Context, task *asynq.Task) error {
	return nil
}

func handleFeedRegenerate(ctx context.Context, task *asynq.Task) error {
	var payload struct {
		TenantID uint64 `json:"tenant_id"`
		FeedID   uint64 `json:"feed_id"`
	}
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return err
	}

	job := model.SyncJob{
		BaseModel: model.BaseModel{TenantID: payload.TenantID},
		Type:      "feed:regenerate",
		Scope:     "feed",
		ObjectID:  &payload.FeedID,
		Status:    "running",
		StartedAt: func() *time.Time { t := time.Now().UTC(); return &t }(),
	}
	database.DB.Create(&job)

	var feedModel model.Feed
	if err := database.DB.Where("id = ? AND tenant_id = ?", payload.FeedID, payload.TenantID).First(&feedModel).Error; err != nil {
		job.Status = "failed"
		job.ErrorMsg = fmt.Sprintf("feed not found: %v", err)
		database.DB.Save(&job)
		return fmt.Errorf("feed not found: %w", err)
	}

	svc := feed.NewFeedService()
	_, err := svc.Generate(ctx, &feedModel)
	if err != nil {
		job.Status = "failed"
		job.ErrorMsg = err.Error()
		database.DB.Save(&job)
		return err
	}

	now := time.Now().UTC()
	job.Status = "success"
	job.FinishedAt = &now
	job.RetryCount = 0
	database.DB.Save(&job)

	return nil
}

func handleRuleEvaluate(ctx context.Context, task *asynq.Task) error {
	var payload struct {
		TenantID uint64 `json:"tenant_id"`
		RuleID   uint64 `json:"rule_id"`
	}
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return err
	}

	job := model.SyncJob{
		BaseModel: model.BaseModel{TenantID: payload.TenantID},
		Type:      "rule:evaluate",
		Scope:     "rule",
		ObjectID:  &payload.RuleID,
		Status:    "running",
		StartedAt: func() *time.Time { t := time.Now().UTC(); return &t }(),
	}
	database.DB.Create(&job)

	var ruleModel model.Rule
	if err := database.DB.Where("id = ? AND tenant_id = ?", payload.RuleID, payload.TenantID).First(&ruleModel).Error; err != nil {
		job.Status = "failed"
		job.ErrorMsg = fmt.Sprintf("rule not found: %v", err)
		database.DB.Save(&job)
		return fmt.Errorf("rule not found: %w", err)
	}

	engine := rule.NewRuleEngine()
	_, err := engine.Evaluate(ctx, &ruleModel)
	if err != nil {
		job.Status = "failed"
		job.ErrorMsg = err.Error()
		database.DB.Save(&job)
		return err
	}

	now := time.Now().UTC()
	job.Status = "success"
	job.FinishedAt = &now
	job.RetryCount = 0
	database.DB.Save(&job)

	return nil
}
