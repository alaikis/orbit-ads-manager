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

type ShopifyAdapter struct {
	config *config.Config
	client *http.Client
}

func NewShopifyAdapter(cfg *config.Config) *ShopifyAdapter {
	return &ShopifyAdapter{config: cfg, client: &http.Client{Timeout: 30 * time.Second}}
}

func (a *ShopifyAdapter) Connect(ctx context.Context, cfg interface{}) (*StoreInfo, error) {
	return nil, fmt.Errorf("use TestConnection instead")
}

func (a *ShopifyAdapter) TestConnection(ctx context.Context) error {
	return fmt.Errorf("not implemented")
}

func (a *ShopifyAdapter) GetStoreInfo(ctx context.Context) (*StoreInfo, error) {
	return nil, fmt.Errorf("not implemented")
}

func (a *ShopifyAdapter) SyncProducts(ctx context.Context, cursor interface{}, limit int) ([]Product, interface{}, error) {
	return nil, nil, fmt.Errorf("not implemented")
}

func (a *ShopifyAdapter) SyncOrders(ctx context.Context, since string, cursor interface{}) ([]Order, interface{}, error) {
	return nil, nil, fmt.Errorf("not implemented")
}

func (a *ShopifyAdapter) FetchInventory(ctx context.Context, productIDs []string) (map[string]int, error) {
	return nil, fmt.Errorf("not implemented")
}

func (a *ShopifyAdapter) SyncProductsFromAPI(storeID uint64, baseURL, accessToken string) ([]Product, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	pageInfo := ""
	var allProducts []Product

	for {
		url := fmt.Sprintf("%s/admin/api/2024-01/products.json?limit=250", baseURL)
		if pageInfo != "" {
			url = fmt.Sprintf("%s/admin/api/2024-01/products.json?limit=250&page_info=%s", baseURL, pageInfo)
		}

		req, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("X-Shopify-Access-Token", accessToken)

		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return nil, fmt.Errorf("shopify API returned status %d: %s", resp.StatusCode, string(body))
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		var result struct {
			Products []struct {
				ID          int64   `json:"id"`
				Title       string  `json:"title"`
				BodyHTML    string  `json:"body_html"`
				Handle      string  `json:"handle"`
				Status      string  `json:"status"`
				Vendor      string  `json:"vendor"`
				ProductType string  `json:"product_type"`
				Tags        string  `json:"tags"`
				Variants    []struct {
					ID           int64   `json:"id"`
					Title        string  `json:"title"`
					SKU          string  `json:"sku"`
					Price        string  `json:"price"`
					CompareAt    string  `json:"compare_at_price"`
					InventoryQty int64   `json:"inventory_quantity"`
					ImageID      *int64  `json:"image_id"`
				} `json:"variants"`
				Images []struct {
					ID    int64  `json:"id"`
					Src   string `json:"src"`
					Width int    `json:"width"`
					Alt   string `json:"alt"`
				} `json:"images"`
				Options []struct {
					Name   string   `json:"name"`
					Values []string `json:"values"`
				} `json:"options"`
			} `json:"products"`
			NextPageInfo string `json:"next_page_info"`
		}

		if err := json.Unmarshal(body, &result); err != nil {
			return nil, err
		}

		for _, p := range result.Products {
			product := Product{
				ID:          fmt.Sprintf("%d", p.ID),
				Title:       p.Title,
				Description: p.BodyHTML,
				Status:      p.Status,
				Brand:       p.Vendor,
				Gtin:        "",
			}

			if len(p.Images) > 0 {
				product.Link = fmt.Sprintf("%s/products/%s", baseURL, p.Handle)
				product.ImageURL = p.Images[0].Src
			}

			if len(p.Variants) > 0 {
				v := p.Variants[0]
				product.PriceCents = parsePrice(v.Price)
				product.SKU = v.SKU
				product.StockQty = int(v.InventoryQty)
			} else {
				product.PriceCents = 0
			}

			for _, v := range p.Variants {
				optionParts := []string{}
				if len(p.Options) > 0 && len(v.Title) > 0 {
					optionParts = append(optionParts, v.Title)
				}
				variant := ProductVariant{
					ID:          fmt.Sprintf("%d", v.ID),
					SKU:         v.SKU,
					OptionString: fmt.Sprintf("%v", optionParts),
					PriceCents:  parsePrice(v.Price),
					InventoryQty: int(v.InventoryQty),
				}
				product.Variants = append(product.Variants, variant)
			}

			allProducts = append(allProducts, product)
		}

		if result.NextPageInfo == "" {
			break
		}
		pageInfo = result.NextPageInfo
	}

	return allProducts, nil
}

func (a *ShopifyAdapter) SyncOrdersFromAPI(storeID uint64, baseURL, accessToken string, since string) ([]Order, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	pageInfo := ""
	var allOrders []Order

	for {
		url := fmt.Sprintf("%s/admin/api/2024-01/orders.json?limit=250&status=any", baseURL)
		if since != "" {
			url = fmt.Sprintf("%s&created_at_min=%s", url, since)
		}
		if pageInfo != "" {
			url = fmt.Sprintf("%s&page_info=%s", url, pageInfo)
		}

		req, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("X-Shopify-Access-Token", accessToken)

		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return nil, fmt.Errorf("shopify API returned status %d: %s", resp.StatusCode, string(body))
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		var result struct {
			Orders []struct {
				ID                int64  `json:"id"`
				Name              string `json:"name"`
				Status            string `json:"financial_status"`
				TotalPrice        string `json:"total_price"`
				Currency          string `json:"currency"`
				CreatedAt         string `json:"created_at"`
				Customer          struct {
					Email string `json:"email"`
				} `json:"customer"`
				LineItems []struct {
					ProductID int64  `json:"product_id"`
					Title     string `json:"title"`
					Quantity  int    `json:"quantity"`
					Price     string `json:"price"`
				} `json:"line_items"`
			} `json:"orders"`
			NextPageInfo string `json:"next_page_info"`
		}

		if err := json.Unmarshal(body, &result); err != nil {
			return nil, err
		}

		for _, o := range result.Orders {
			order := Order{
				ID:            fmt.Sprintf("%d", o.ID),
				OrderNumber:   o.Name,
				Status:        o.Status,
				TotalCents:    parsePrice(o.TotalPrice),
				Currency:      o.Currency,
				CustomerEmail: o.Customer.Email,
				PlacedAt:      o.CreatedAt,
			}

			items := make([]map[string]interface{}, 0, len(o.LineItems))
			for _, item := range o.LineItems {
				items = append(items, map[string]interface{}{
					"product_id": item.ProductID,
					"name":       item.Title,
					"quantity":   item.Quantity,
					"total":      item.Price,
				})
			}
			order.Items = items

			allOrders = append(allOrders, order)
		}

		if result.NextPageInfo == "" {
			break
		}
		pageInfo = result.NextPageInfo
	}

	return allOrders, nil
}
