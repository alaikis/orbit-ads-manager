package store

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"orbit/apps/api/config"
)

type WooCommerceAdapter struct {
	config *config.Config
	client *http.Client
}

func NewWooCommerceAdapter(cfg *config.Config) *WooCommerceAdapter {
	return &WooCommerceAdapter{config: cfg, client: &http.Client{Timeout: 30 * time.Second}}
}

func (a *WooCommerceAdapter) Connect(ctx context.Context, cfg interface{}) (*StoreInfo, error) {
	return nil, fmt.Errorf("use TestConnection instead")
}

func (a *WooCommerceAdapter) TestConnection(ctx context.Context) error {
	return fmt.Errorf("not implemented")
}

func (a *WooCommerceAdapter) GetStoreInfo(ctx context.Context) (*StoreInfo, error) {
	return nil, fmt.Errorf("not implemented")
}

func (a *WooCommerceAdapter) SyncProducts(ctx context.Context, cursor interface{}, limit int) ([]Product, interface{}, error) {
	return nil, nil, fmt.Errorf("not implemented")
}

func (a *WooCommerceAdapter) SyncOrders(ctx context.Context, since string, cursor interface{}) ([]Order, interface{}, error) {
	return nil, nil, fmt.Errorf("not implemented")
}

func (a *WooCommerceAdapter) FetchInventory(ctx context.Context, productIDs []string) (map[string]int, error) {
	return nil, fmt.Errorf("not implemented")
}

func (a *WooCommerceAdapter) SyncProductsFromAPI(storeID uint64, baseURL, apiKey, apiSecret string) ([]Product, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	page := 1
	var allProducts []Product

	for {
		url := fmt.Sprintf("%s/wp-json/wc/v3/products?per_page=100&page=%d", baseURL, page)
		req, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
		if err != nil {
			return nil, err
		}
		req.SetBasicAuth(apiKey, apiSecret)

		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return nil, fmt.Errorf("woocommerce API returned status %d: %s", resp.StatusCode, string(body))
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		var raw []map[string]interface{}
		if err := json.Unmarshal(body, &raw); err != nil {
			return nil, err
		}

		if len(raw) == 0 {
			break
		}

		for _, p := range raw {
			product := Product{
				ID:          fmt.Sprintf("%v", p["id"]),
				Title:       fmt.Sprintf("%v", p["name"]),
				Description: fmt.Sprintf("%v", p["description"]),
				Status:      fmt.Sprintf("%v", p["status"]),
				Brand:       fmt.Sprintf("%v", p["brand"]),
				Gtin:        fmt.Sprintf("%v", p["gtin"]),
			}

			if link, ok := p["permalink"].(string); ok {
				product.Link = link
			}

			if images, ok := p["images"].([]interface{}); ok && len(images) > 0 {
				if imgMap, ok := images[0].(map[string]interface{}); ok {
					if src, ok := imgMap["src"].(string); ok {
						product.ImageURL = src
					}
				}
			}

			if priceStr, ok := p["price"].(string); ok && priceStr != "" {
				product.PriceCents = parsePrice(priceStr)
			} else if regularPriceStr, ok := p["regular_price"].(string); ok && regularPriceStr != "" {
				product.PriceCents = parsePrice(regularPriceStr)
			}

			if sku, ok := p["sku"].(string); ok {
				product.SKU = sku
			}

			if stockQty, ok := p["stock_quantity"].(float64); ok {
				product.StockQty = int(stockQty)
			}

			if variants, ok := p["variations"].([]interface{}); ok && len(variants) > 0 {
				for _, v := range variants {
					if vMap, ok := v.(map[string]interface{}); ok {
						variant := ProductVariant{
							ID:          fmt.Sprintf("%v", vMap["id"]),
							SKU:         fmt.Sprintf("%v", vMap["sku"]),
							OptionString: fmt.Sprintf("%v", vMap["attributes"]),
							PriceCents:  product.PriceCents,
						}
						if stockQty, ok := vMap["stock_quantity"].(float64); ok {
							variant.InventoryQty = int(stockQty)
						}
						product.Variants = append(product.Variants, variant)
					}
				}
			}

			allProducts = append(allProducts, product)
		}

		if len(raw) < 100 {
			break
		}
		page++
	}

	return allProducts, nil
}

func (a *WooCommerceAdapter) SyncOrdersFromAPI(storeID uint64, baseURL, apiKey, apiSecret string, since string) ([]Order, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	page := 1
	var allOrders []Order

	for {
		url := fmt.Sprintf("%s/wp-json/wc/v3/orders?per_page=100&page=%d&after=%s", baseURL, page, since)
		req, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
		if err != nil {
			return nil, err
		}
		req.SetBasicAuth(apiKey, apiSecret)

		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return nil, fmt.Errorf("woocommerce API returned status %d: %s", resp.StatusCode, string(body))
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		var raw []map[string]interface{}
		if err := json.Unmarshal(body, &raw); err != nil {
			return nil, err
		}

		if len(raw) == 0 {
			break
		}

		for _, o := range raw {
			order := Order{
				ID:            fmt.Sprintf("%v", o["id"]),
				OrderNumber:   fmt.Sprintf("%v", o["number"]),
				Status:        fmt.Sprintf("%v", o["status"]),
				CustomerEmail: fmt.Sprintf("%v", o["billing_email"]),
			}

			if totalStr, ok := o["total"].(string); ok {
				order.TotalCents = parsePrice(totalStr)
			}
			if currency, ok := o["currency"].(string); ok {
				order.Currency = currency
			}
			if dateCreated, ok := o["date_created"].(string); ok && dateCreated != "" {
				order.PlacedAt = dateCreated
			}

			if lineItems, ok := o["line_items"].([]interface{}); ok {
				items := make([]map[string]interface{}, 0, len(lineItems))
				for _, item := range lineItems {
					if itemMap, ok := item.(map[string]interface{}); ok {
						items = append(items, map[string]interface{}{
							"product_id": itemMap["product_id"],
							"name":       itemMap["name"],
							"quantity":   itemMap["quantity"],
							"total":      itemMap["total"],
						})
					}
				}
				order.Items = items
			}

			allOrders = append(allOrders, order)
		}

		if len(raw) < 100 {
			break
		}
		page++
	}

	return allOrders, nil
}
