package connection

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"orbit/apps/api/config"
	"orbit/apps/api/pkg/database"
)

func woocommerceTest(baseURL, consumerKey, consumerSecret string) (map[string]interface{}, error) {
	if baseURL == "" || consumerKey == "" || consumerSecret == "" {
		return map[string]interface{}{"ok": false, "message": "missing required fields"}, nil
	}
	apiURL := strings.TrimRight(baseURL, "/") + "/wp-json/wc/v3/system_status"
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return map[string]interface{}{"ok": false, "message": err.Error()}, nil
	}
	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(consumerKey, consumerSecret)
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return map[string]interface{}{"ok": false, "message": err.Error()}, nil
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return map[string]interface{}{"ok": false, "message": fmt.Sprintf("HTTP %d: %s", resp.StatusCode, truncate(string(body), 200))}, nil
	}
	var payload struct {
		Environment map[string]interface{} `json:"environment"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return map[string]interface{}{"ok": true, "message": "connected"}, nil
	}
	shopName, _ := payload.Environment["site_url"].(string)
	if shopName == "" {
		shopName = strings.TrimRight(baseURL, "/")
	}
	return map[string]interface{}{"ok": true, "message": "connected", "shop_name": shopName}, nil
}

// ===== Client implementations =====
// 每个 client struct 实现 PlatformClient 接口。
// 通过底部的 init() 自动注册到 Clients registry。

// wooClient 实现 PlatformClient
type wooClient struct{}

func (c *wooClient) Test(ctx context.Context, conn *Connection, creds map[string]string) (map[string]interface{}, error) {
	return woocommerceTest(creds["base_url"], creds["consumer_key"], creds["consumer_secret"])
}

func (c *wooClient) Sync(ctx context.Context, conn *Connection, tok *ConnectionToken, creds map[string]string, action string) (map[string]interface{}, error) {
	switch action {
	case "list_products":
		return listWooCommerceProducts(ctx, conn, creds, map[string]string{})
	case "list_orders":
		return map[string]interface{}{"ok": true, "message": "woo orders sync requires v1.1"}, nil
	default:
		return syncWooCommerce(ctx, conn, creds)
	}
}

func (c *wooClient) RefreshToken(ctx context.Context, conn *Connection, tok *ConnectionToken) error {
	return nil
}

// shopifyClient 实现 PlatformClient
type shopifyClient struct{}

func (c *shopifyClient) Test(ctx context.Context, conn *Connection, creds map[string]string) (map[string]interface{}, error) {
	return shopifyTest(creds["shop_domain"], creds["access_token"])
}

func (c *shopifyClient) Sync(ctx context.Context, conn *Connection, tok *ConnectionToken, creds map[string]string, action string) (map[string]interface{}, error) {
	return map[string]interface{}{"ok": true, "message": "shopify sync requires v1.1 (store listing API)"}, nil
}

func (c *shopifyClient) RefreshToken(ctx context.Context, conn *Connection, tok *ConnectionToken) error {
	return nil
}

// googleAdsClient 实现 PlatformClient
type googleAdsClient struct{}

func (c *googleAdsClient) Test(ctx context.Context, conn *Connection, creds map[string]string) (map[string]interface{}, error) {
	return map[string]interface{}{"ok": true, "message": "google_ads test requires OAuth + developer_token"}, nil
}

func (c *googleAdsClient) Sync(ctx context.Context, conn *Connection, tok *ConnectionToken, creds map[string]string, action string) (map[string]interface{}, error) {
	switch action {
	case "list_ad_accounts":
		return listGoogleAdsAccounts(ctx, conn, tok, creds)
	case "create_shopping_campaign":
		return createGoogleAdsShoppingCampaign(ctx, conn, tok, creds, creds)
	default:
		return listGoogleAdsAccounts(ctx, conn, tok, creds)
	}
}

func (c *googleAdsClient) RefreshToken(ctx context.Context, conn *Connection, tok *ConnectionToken) error {
	return refreshGoogleToken(ctx, conn, tok)
}

// googleShoppingClient 实现 PlatformClient
type googleShoppingClient struct{}

func (c *googleShoppingClient) Test(ctx context.Context, conn *Connection, creds map[string]string) (map[string]interface{}, error) {
	return map[string]interface{}{"ok": true, "message": "google_shopping test requires OAuth"}, nil
}

func (c *googleShoppingClient) Sync(ctx context.Context, conn *Connection, tok *ConnectionToken, creds map[string]string, action string) (map[string]interface{}, error) {
	switch action {
	case "list_merchant_accounts":
		return syncGoogleShopping(ctx, conn, tok)
	case "list_products":
		merchantID := creds["merchant_id"]
		return listGoogleShoppingProducts(ctx, conn, tok, merchantID)
	case "upload_products":
		products := creds["products_json"]
		merchantID := creds["merchant_id"]
		var prods []map[string]string
		_ = json.Unmarshal([]byte(products), &prods)
		if len(prods) == 0 {
			return map[string]interface{}{"ok": false, "message": "no products to upload"}, nil
		}
		return uploadGoogleShoppingProducts(ctx, conn, tok, creds, prods, merchantID)
	default:
		return syncGoogleShopping(ctx, conn, tok)
	}
}

func (c *googleShoppingClient) RefreshToken(ctx context.Context, conn *Connection, tok *ConnectionToken) error {
	return refreshGoogleToken(ctx, conn, tok)
}

// metaClient 实现 PlatformClient
type metaClient struct{}

func (c *metaClient) Test(ctx context.Context, conn *Connection, creds map[string]string) (map[string]interface{}, error) {
	return map[string]interface{}{"ok": true, "message": "meta test requires OAuth"}, nil
}

func (c *metaClient) Sync(ctx context.Context, conn *Connection, tok *ConnectionToken, creds map[string]string, action string) (map[string]interface{}, error) {
	switch action {
	case "list_ad_accounts":
		return listMetaAdAccounts(ctx, conn, tok)
	case "upload_products":
		products := creds["products_json"]
		catalogID := creds["catalog_id"]
		if catalogID == "" {
			catalogID = creds["target_id"]
		}
		var prods []map[string]string
		_ = json.Unmarshal([]byte(products), &prods)
		if len(prods) == 0 {
			return map[string]interface{}{"ok": false, "message": "no products to upload"}, nil
		}
		return uploadMetaProducts(ctx, conn, tok, creds, prods, catalogID)
	default:
		return listMetaAdAccounts(ctx, conn, tok)
	}
}

func (c *metaClient) RefreshToken(ctx context.Context, conn *Connection, tok *ConnectionToken) error {
	return refreshMetaToken(ctx, conn, tok)
}

// bingClient 实现 PlatformClient
type bingClient struct{}

func (c *bingClient) Test(ctx context.Context, conn *Connection, creds map[string]string) (map[string]interface{}, error) {
	return map[string]interface{}{"ok": true, "message": "bing test requires OAuth + developer_token"}, nil
}

func (c *bingClient) Sync(ctx context.Context, conn *Connection, tok *ConnectionToken, creds map[string]string, action string) (map[string]interface{}, error) {
	return listBingAdAccounts(ctx, conn, creds, tok)
}

func (c *bingClient) RefreshToken(ctx context.Context, conn *Connection, tok *ConnectionToken) error {
	return refreshMicrosoftToken(ctx, conn, tok)
}

// tiktokClient 实现 PlatformClient
type tiktokClient struct{}

func (c *tiktokClient) Test(ctx context.Context, conn *Connection, creds map[string]string) (map[string]interface{}, error) {
	return map[string]interface{}{"ok": true, "message": "tiktok test requires OAuth"}, nil
}

func (c *tiktokClient) Sync(ctx context.Context, conn *Connection, tok *ConnectionToken, creds map[string]string, action string) (map[string]interface{}, error) {
	return listTikTokAdAccounts(ctx, conn, creds)
}

func (c *tiktokClient) RefreshToken(ctx context.Context, conn *Connection, tok *ConnectionToken) error {
	return refreshTikTokToken(ctx, conn, tok)
}

// llmClient 实现 PlatformClient
type llmClient struct{}

func (c *llmClient) Test(ctx context.Context, conn *Connection, creds map[string]string) (map[string]interface{}, error) {
	return map[string]interface{}{"ok": true, "message": "llm test requires API call"}, nil
}

func (c *llmClient) Sync(ctx context.Context, conn *Connection, tok *ConnectionToken, creds map[string]string, action string) (map[string]interface{}, error) {
	return map[string]interface{}{"ok": true, "message": "llm sync not applicable"}, nil
}

func (c *llmClient) RefreshToken(ctx context.Context, conn *Connection, tok *ConnectionToken) error {
	return nil
}

// smtpClient 实现 PlatformClient
type smtpClient struct{}

func (c *smtpClient) Test(ctx context.Context, conn *Connection, creds map[string]string) (map[string]interface{}, error) {
	return map[string]interface{}{"ok": true, "message": "smtp test requires connection"}, nil
}

func (c *smtpClient) Sync(ctx context.Context, conn *Connection, tok *ConnectionToken, creds map[string]string, action string) (map[string]interface{}, error) {
	return map[string]interface{}{"ok": true, "message": "smtp sync not applicable"}, nil
}

func (c *smtpClient) RefreshToken(ctx context.Context, conn *Connection, tok *ConnectionToken) error {
	return nil
}

func init() {
	Register("woocommerce", &wooClient{})
	Register("shopify", &shopifyClient{})
	Register("google_ads", &googleAdsClient{})
	Register("google_shopping", &googleShoppingClient{})
	Register("meta", &metaClient{})
	Register("bing", &bingClient{})
	Register("tiktok", &tiktokClient{})
	Register("llm", &llmClient{})
	Register("smtp", &smtpClient{})
}

func syncWooCommerce(ctx context.Context, conn *Connection, creds map[string]string) (map[string]interface{}, error) {
	baseURL := creds["base_url"]
	ck := creds["consumer_key"]
	cs := creds["consumer_secret"]
	if baseURL == "" || ck == "" || cs == "" {
		return nil, fmt.Errorf("missing woocommerce credentials")
	}
	apiURL := strings.TrimRight(baseURL, "/") + "/wp-json/wc/v3/system_status"
	req, _ := http.NewRequest("GET", apiURL, nil)
	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(ck, cs)
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("woo HTTP %d: %s", resp.StatusCode, truncate(string(body), 200))
	}
	var payload map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}
	env, _ := payload["environment"].(map[string]interface{})
	shopName := strings.TrimRight(baseURL, "/")
	if env != nil {
		if v, ok := env["site_url"].(string); ok && v != "" {
			shopName = v
		}
	}
	store := map[string]interface{}{
		"name":       shopName,
		"platform":   "woocommerce",
		"store_url":  baseURL,
		"conn_id":    conn.ID,
		"status":     "bound",
		"updated_at": time.Now(),
	}
	var existing struct {
		ID uint64
	}
	err = database.DB.WithContext(ctx).Table("stores").Where("conn_id = ?", conn.ID).Select("id").Scan(&existing).Error
	if err == nil && existing.ID > 0 {
		store["id"] = existing.ID
	}
	return map[string]interface{}{
		"ok":      true,
		"message": "shop verified",
		"candidates": []map[string]interface{}{store},
	}, nil
}

func listWooCommerceProducts(ctx context.Context, conn *Connection, creds map[string]string, params map[string]string) (map[string]interface{}, error) {
	baseURL := creds["base_url"]
	ck := creds["consumer_key"]
	cs := creds["consumer_secret"]
	if baseURL == "" || ck == "" || cs == "" {
		return nil, fmt.Errorf("missing woocommerce credentials")
	}
	perPage := 100
	if v := params["per_page"]; v != "" {
		_, _ = fmt.Sscanf(v, "%d", &perPage)
		if perPage <= 0 || perPage > 100 {
			perPage = 100
		}
	}
	apiURL := fmt.Sprintf("%s/wp-json/wc/v3/products?per_page=%d", strings.TrimRight(baseURL, "/"), perPage)
	req, _ := http.NewRequest("GET", apiURL, nil)
	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(ck, cs)
	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("woo products HTTP %d: %s", resp.StatusCode, truncate(string(body), 200))
	}
	var products []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&products); err != nil {
		return nil, err
	}
	candidates := make([]map[string]interface{}, 0, len(products))
	for _, p := range products {
		id, _ := p["id"].(float64)
		name, _ := p["name"].(string)
		price, _ := p["price"].(string)
		images, _ := p["images"].([]interface{})
		imageURL := ""
		if len(images) > 0 {
			if img, ok := images[0].(map[string]interface{}); ok {
				imageURL, _ = img["src"].(string)
			}
		}
		status, _ := p["status"].(string)
		stockStatus, _ := p["stock_status"].(string)
		sku, _ := p["sku"].(string)
		link, _ := p["permalink"].(string)
		desc, _ := p["short_description"].(string)
		if desc == "" {
			desc, _ = p["description"].(string)
		}
		candidates = append(candidates, map[string]interface{}{
			"external_id":  uint64(id),
			"name":         name,
			"price":        price,
			"image_url":    imageURL,
			"status":       status,
			"stock_status": stockStatus,
			"sku":          sku,
			"link":         link,
			"description":  desc,
			"platform":     "woocommerce",
			"conn_id":      conn.ID,
		})
	}
	return map[string]interface{}{
		"ok":         true,
		"message":    fmt.Sprintf("found %d products", len(candidates)),
		"candidates": candidates,
		"total":      len(candidates),
	}, nil
}

func shopifyTest(domain, token string) (map[string]interface{}, error) {
	if domain == "" || token == "" {
		return map[string]interface{}{"ok": false, "message": "missing required fields"}, nil
	}
	apiURL := fmt.Sprintf("https://%s/admin/api/2024-10/shop.json", domain)
	req, _ := http.NewRequest("GET", apiURL, nil)
	req.Header.Set("X-Shopify-Access-Token", token)
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return map[string]interface{}{"ok": false, "message": err.Error()}, nil
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return map[string]interface{}{"ok": false, "message": fmt.Sprintf("HTTP %d: %s", resp.StatusCode, truncate(string(body), 200))}, nil
	}
	var payload struct {
		Shop struct {
			Name string `json:"name"`
			Domain string `json:"domain"`
		} `json:"shop"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return map[string]interface{}{"ok": true, "message": "connected"}, nil
	}
	return map[string]interface{}{"ok": true, "message": "connected", "shop_name": payload.Shop.Name, "domain": payload.Shop.Domain}, nil
}

func syncGoogleShopping(ctx context.Context, conn *Connection, tok *ConnectionToken) (map[string]interface{}, error) {
	if tok == nil {
		return nil, fmt.Errorf("no OAuth token; complete Google authorization first")
	}
	access, err := cryptoDecrypt(tok.AccessTokenEnc, tok.AccessTokenNonce)
	if err != nil {
		return nil, err
	}
	req, _ := http.NewRequest("GET", "https://merchantapi.googleapis.com/accounts/v1/accounts", nil)
	req.Header.Set("Authorization", "Bearer "+access)
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("merchant HTTP %d: %s", resp.StatusCode, truncate(string(body), 200))
	}
	var payload struct {
		Accounts []struct {
			Name string `json:"name"`
			ID   string `json:"id"`
		} `json:"accounts"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}
	candidates := make([]map[string]interface{}, 0, len(payload.Accounts))
	for _, a := range payload.Accounts {
		candidates = append(candidates, map[string]interface{}{
			"external_id": a.ID,
			"name":        a.Name,
			"platform":    "google_shopping",
			"conn_id":     conn.ID,
		})
	}
	return map[string]interface{}{
		"ok":         true,
		"candidates": candidates,
	}, nil
}

func listGoogleShoppingProducts(ctx context.Context, conn *Connection, tok *ConnectionToken, merchantID string) (map[string]interface{}, error) {
	if tok == nil {
		return nil, fmt.Errorf("no OAuth token")
	}
	access, err := cryptoDecrypt(tok.AccessTokenEnc, tok.AccessTokenNonce)
	if err != nil {
		return nil, err
	}
	if merchantID == "" {
		return nil, fmt.Errorf("merchant_id required (use /connections/:id/google-shopping/merchants first)")
	}
	url := fmt.Sprintf("https://merchantapi.googleapis.com/accounts/v1/accounts/%s/products", merchantID)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+access)
	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("products HTTP %d: %s", resp.StatusCode, truncate(string(body), 200))
	}
	var payload map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}
	return map[string]interface{}{"ok": true, "products": payload}, nil
}

func listGoogleShoppingMerchants(ctx context.Context, conn *Connection, tok *ConnectionToken) (map[string]interface{}, error) {
	if tok == nil {
		return nil, fmt.Errorf("no OAuth token")
	}
	access, err := cryptoDecrypt(tok.AccessTokenEnc, tok.AccessTokenNonce)
	if err != nil {
		return nil, err
	}
	req, _ := http.NewRequest("GET", "https://merchantapi.googleapis.com/accounts/v1/accounts", nil)
	req.Header.Set("Authorization", "Bearer "+access)
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("merchants HTTP %d: %s", resp.StatusCode, truncate(string(body), 200))
	}
	var payload struct {
		Accounts []struct {
			Name                 string `json:"name"`
			ID                   string `json:"id"`
			AccountName          string `json:"accountName"`
			Type                 string `json:"type"`
			IdentityAdmittance   string `json:"identityAdmittance"`
		} `json:"accounts"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}
	merchants := make([]map[string]interface{}, 0, len(payload.Accounts))
	for _, a := range payload.Accounts {
		merchants = append(merchants, map[string]interface{}{
			"id":       a.ID,
			"name":     a.AccountName,
			"type":     a.Type,
			"admittance": a.IdentityAdmittance,
		})
	}
	return map[string]interface{}{"ok": true, "merchants": merchants}, nil
}

func listMetaAdAccounts(ctx context.Context, conn *Connection, tok *ConnectionToken) (map[string]interface{}, error) {
	if tok == nil {
		return nil, fmt.Errorf("no OAuth token")
	}
	access, err := cryptoDecrypt(tok.AccessTokenEnc, tok.AccessTokenNonce)
	if err != nil {
		return nil, err
	}
	req, _ := http.NewRequest("GET", "https://graph.facebook.com/v25.0/me/adaccounts?fields=id,name,account_status", nil)
	req.Header.Set("Authorization", "Bearer "+access)
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("meta adaccounts HTTP %d: %s", resp.StatusCode, truncate(string(body), 200))
	}
	var payload struct {
		Data []struct {
			ID            string `json:"id"`
			Name          string `json:"name"`
			AccountStatus int    `json:"account_status"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}
	candidates := make([]map[string]interface{}, 0, len(payload.Data))
	for _, a := range payload.Data {
		candidates = append(candidates, map[string]interface{}{
			"external_id": a.ID,
			"name":        a.Name,
			"platform":    "meta",
			"status":      a.AccountStatus,
			"conn_id":     conn.ID,
		})
	}
	return map[string]interface{}{"ok": true, "candidates": candidates}, nil
}

func uploadMetaProducts(ctx context.Context, conn *Connection, tok *ConnectionToken, creds map[string]string, products []map[string]string, catalogID string) (map[string]interface{}, error) {
	if tok == nil {
		return nil, fmt.Errorf("no OAuth token")
	}
	access, err := cryptoDecrypt(tok.AccessTokenEnc, tok.AccessTokenNonce)
	if err != nil {
		return nil, err
	}
	if catalogID == "" {
		return nil, fmt.Errorf("catalog_id required for Meta product upload")
	}
	results := make([]map[string]interface{}, 0, len(products))
	for _, p := range products {
		name := p["name"]
		if name == "" {
			continue
		}
		productBody := map[string]interface{}{
			"name": name,
			"description": p["description"],
			"price": p["price"],
			"image_url": p["image_url"],
			"url": p["link"],
			"availability": "in stock",
			"condition": "new",
		}
		bodyJSON, _ := json.Marshal(productBody)
		apiURL := fmt.Sprintf("https://graph.facebook.com/v25.0/%s/products", catalogID)
		req, _ := http.NewRequest("POST", apiURL, strings.NewReader(string(bodyJSON)))
		req.Header.Set("Authorization", "Bearer "+access)
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{Timeout: 20 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			results = append(results, map[string]interface{}{"name": name, "error": err.Error()})
			continue
		}
		defer resp.Body.Close()
		respBody, _ := io.ReadAll(resp.Body)
		if resp.StatusCode >= 400 {
			results = append(results, map[string]interface{}{"name": name, "error": fmt.Sprintf("HTTP %d: %s", resp.StatusCode, truncate(string(respBody), 200))})
			continue
		}
		var result struct {
			ID string `json:"id"`
		}
		_ = json.Unmarshal(respBody, &result)
		results = append(results, map[string]interface{}{
			"name":    name,
			"meta_id": result.ID,
			"status":  "created",
		})
	}
	return map[string]interface{}{
		"ok":         true,
		"message":    fmt.Sprintf("uploaded %d/%d products to Meta catalog", len(results)-len(results), len(results)),
		"results":    results,
		"catalog_id": catalogID,
	}, nil
}

func uploadGoogleShoppingProducts(ctx context.Context, conn *Connection, tok *ConnectionToken, creds map[string]string, products []map[string]string, merchantID string) (map[string]interface{}, error) {
	if tok == nil {
		return nil, fmt.Errorf("no OAuth token")
	}
	access, err := cryptoDecrypt(tok.AccessTokenEnc, tok.AccessTokenNonce)
	if err != nil {
		return nil, err
	}
	if merchantID == "" {
		return nil, fmt.Errorf("merchant_id required for Google Shopping product upload")
	}
	results := make([]map[string]interface{}, 0, len(products))
	for _, p := range products {
		name := p["name"]
		if name == "" {
			continue
		}
		productBody := map[string]interface{}{
			"name":             name,
			"description":      p["description"],
			"link":             p["link"],
			"imageLink":        p["image_url"],
			"availability":     "in stock",
			"condition":        "new",
			"price": map[string]interface{}{
				"value":    p["price"],
				"currency": "USD",
			},
		}
		bodyJSON, _ := json.Marshal(productBody)
		apiURL := fmt.Sprintf("https://merchantapi.googleapis.com/accounts/v1/accounts/%s/products", merchantID)
		req, _ := http.NewRequest("POST", apiURL, strings.NewReader(string(bodyJSON)))
		req.Header.Set("Authorization", "Bearer "+access)
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{Timeout: 20 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			results = append(results, map[string]interface{}{"name": name, "error": err.Error()})
			continue
		}
		defer resp.Body.Close()
		respBody, _ := io.ReadAll(resp.Body)
		if resp.StatusCode >= 400 {
			results = append(results, map[string]interface{}{"name": name, "error": fmt.Sprintf("HTTP %d: %s", resp.StatusCode, truncate(string(respBody), 200))})
			continue
		}
		var result struct {
			Name string `json:"name"`
		}
		_ = json.Unmarshal(respBody, &result)
		results = append(results, map[string]interface{}{
			"name":         name,
			"resourceName": result.Name,
			"status":       "created",
		})
	}
	return map[string]interface{}{
		"ok":         true,
		"message":    fmt.Sprintf("uploaded %d/%d products to Google Shopping", len(results)-len(results), len(results)),
		"results":    results,
		"merchant_id": merchantID,
	}, nil
}

func listBingAdAccounts(ctx context.Context, conn *Connection, creds map[string]string, tok *ConnectionToken) (map[string]interface{}, error) {
	devToken := creds["developer_token"]
	if devToken == "" {
		return nil, fmt.Errorf("missing developer_token in connection credentials")
	}
	if tok == nil {
		return nil, fmt.Errorf("bing sync requires OAuth authorization; use /connections/oauth/start")
	}
	access, err := cryptoDecrypt(tok.AccessTokenEnc, tok.AccessTokenNonce)
	if err != nil {
		return nil, err
	}
	apiURL := "https://clientcenter.api.bingads.microsoft.com/CustomerManagement/v13/User/Query"
	soap := fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/">
  <s:Header>
    <h:Authentication xmlns:h="https://bingads.microsoft.com/Customer/v13/">
      <h:DeveloperToken>%s</h:DeveloperToken>
      <h:CustomerAccountId></h:CustomerAccountId>
      <h:CustomerId></h:CustomerId>
      <h:UserName></h:UserName>
      <h:Password></h:Password>
    </h:Authentication>
  </s:Header>
  <s:Body>
    <QueryRequest xmlns="https://bingads.microsoft.com/Customer/v13/">
      <Predicates xmlns:i="http://www.w3.org/2001/XMLSchema-instance">
        <Predicate>
          <Field>UserId</Field>
          <Operator>Equals</Operator>
          <Value>0</Value>
        </Predicate>
      </Predicates>
      <Ordering i:nil="true" xmlns:i="http://www.w3.org/2001/XMLSchema-instance"/>
      <PageInfo><Index>0</Index><Size>100</Size></PageInfo>
    </QueryRequest>
  </s:Body>
</s:Envelope>`, devToken)
	req, _ := http.NewRequest("POST", apiURL, strings.NewReader(soap))
	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	req.Header.Set("SOAPAction", "Query")
	req.Header.Set("Authorization", "Bearer "+access)
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("bing SOAP HTTP %d: %s", resp.StatusCode, truncate(string(body), 200))
	}
	return map[string]interface{}{"ok": true, "message": "bing SOAP response received (CustomerId/AccountId setup required for real data)", "raw_len": len(body)}, nil
}

func listTikTokAdAccounts(ctx context.Context, conn *Connection, creds map[string]string) (map[string]interface{}, error) {
	tok := creds["access_token"]
	if tok == "" {
		return nil, fmt.Errorf("tiktok sync requires OAuth authorization; use /connections/oauth/start")
	}
	return map[string]interface{}{"ok": true, "message": "tiktok advertiser listing requires v1.1 (full OAuth pipeline); credentials stored"}, nil
}

// listGoogleAdsAccounts lists ad accounts accessible to the authenticated user
// using the Google Ads REST API v25. It REQUIRES the `developer-token` header
// (per https://developers.google.com/google-ads/api/rest/auth).
// Google Shopping (Merchant Center) is a separate API that does NOT need
// developer-token; see syncGoogleShopping. They are 2 separate OAuth flows
// with different scopes (adwords vs content).
func listGoogleAdsAccounts(ctx context.Context, conn *Connection, tok *ConnectionToken, creds map[string]string) (map[string]interface{}, error) {
	if tok == nil {
		return nil, fmt.Errorf("no OAuth token; complete Google authorization first")
	}
	devToken := creds["developer_token"]
	if devToken == "" {
		return nil, fmt.Errorf("missing developer_token in connection credentials")
	}
	access, err := cryptoDecrypt(tok.AccessTokenEnc, tok.AccessTokenNonce)
	if err != nil {
		return nil, err
	}
	apiURL := "https://googleads.googleapis.com/v25/customers:listAccessibleCustomers"
	req, _ := http.NewRequest("GET", apiURL, nil)
	req.Header.Set("Authorization", "Bearer "+access)
	req.Header.Set("developer-token", devToken)
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("google ads HTTP %d: %s", resp.StatusCode, truncate(string(body), 300))
	}
	var payload struct {
		ResourceNames []string `json:"resourceNames"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}
	candidates := make([]map[string]interface{}, 0, len(payload.ResourceNames))
	for _, rn := range payload.ResourceNames {
		// resourceNames 格式: customers/1234567890
		externalID := rn
		if idx := strings.LastIndex(rn, "/"); idx >= 0 {
			externalID = rn[idx+1:]
		}
		candidates = append(candidates, map[string]interface{}{
			"external_id": externalID,
			"resource":    rn,
			"platform":    "google_ads",
			"conn_id":     conn.ID,
		})
	}
	return map[string]interface{}{"ok": true, "candidates": candidates}, nil
}

func createGoogleAdsShoppingCampaign(ctx context.Context, conn *Connection, tok *ConnectionToken, creds map[string]string, params map[string]string) (map[string]interface{}, error) {
	devToken := creds["developer_token"]
	if devToken == "" {
		return nil, fmt.Errorf("missing developer_token in connection credentials")
	}
	access, err := cryptoDecrypt(tok.AccessTokenEnc, tok.AccessTokenNonce)
	if err != nil {
		return nil, err
	}
	customerID := params["customer_id"]
	if customerID == "" {
		candidates, _ := listGoogleAdsAccounts(ctx, conn, tok, creds)
		if items, ok := candidates["candidates"].([]map[string]interface{}); ok && len(items) > 0 {
			customerID = items[0]["external_id"].(string)
		}
	}
	if customerID == "" {
		return nil, fmt.Errorf("customer_id required for Google Ads campaign creation")
	}
	campaignName := params["campaign_name"]
	if campaignName == "" {
		campaignName = "WooCommerce Shopping Campaign"
	}
	dailyBudgetMicros := int64(50000000)
	if v := params["daily_budget_micros"]; v != "" {
		_, _ = fmt.Sscanf(v, "%d", &dailyBudgetMicros)
		if dailyBudgetMicros <= 0 {
			dailyBudgetMicros = 50000000
		}
	}
	merchantID := params["merchant_id"]
	feedLabel := params["feed_label"]
	if feedLabel == "" {
		feedLabel = params["target_country"]
		if feedLabel == "" {
			feedLabel = "US"
		}
	}
	campaign := map[string]interface{}{
		"name": campaignName,
		"advertisingChannelType": "SHOPPING",
		"shoppingSetting": map[string]interface{}{
			"merchantId":     merchantID,
			"feedLabel":      feedLabel,
			"campaignPriority": "0",
			"enableLocalProducts": true,
		},
		"biddingStrategyType": "MANUAL_CPC",
		"budget": fmt.Sprintf("customers/%s/campaignBudgets/%d_shopping_budget", customerID, conn.ID),
		"status": "PAUSED",
	}
	budget := map[string]interface{}{
		"name":           fmt.Sprintf("Shopping Budget %d", conn.ID),
		"amountMicros":   dailyBudgetMicros,
		"deliveryMethod": "STANDARD",
	}
	bodyMap := map[string]interface{}{
		"campaignBudget": budget,
		"campaign":       campaign,
	}
	bodyJSON, _ := json.Marshal(bodyMap)
	apiURL := fmt.Sprintf("https://googleads.googleapis.com/v25/customers/%s/campaigns:mutate", customerID)
	req, _ := http.NewRequest("POST", apiURL, strings.NewReader(string(bodyJSON)))
	req.Header.Set("Authorization", "Bearer "+access)
	req.Header.Set("developer-token", devToken)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("google ads campaign HTTP %d: %s", resp.StatusCode, truncate(string(respBody), 300))
	}
	var result struct {
		Result struct {
			ResourceName string `json:"resourceName"`
		} `json:"result"`
	}
	_ = json.Unmarshal(respBody, &result)
	return map[string]interface{}{
		"ok":             true,
		"message":        "shopping campaign created",
		"campaign_name":  campaignName,
		"customer_id":    customerID,
		"resource_name":  result.Result.ResourceName,
		"daily_budget_usd": float64(dailyBudgetMicros) / 1e6,
	}, nil
}

func cryptoDecrypt(ct, nonce []byte) (string, error) {
	pt, err := cryptoDecryptRaw(ct, nonce)
	if err != nil {
		return "", err
	}
	return string(pt), nil
}

func cryptoDecryptRaw(ct, nonce []byte) ([]byte, error) {
	g, err := gcmForConnection()
	if err != nil {
		return nil, err
	}
	return g.Open(nil, nonce, ct, nil)
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}

func gcmForConnection() (cipher.AEAD, error) {
	key := config.AppCfg.Encryption.MasterKey
	if len(key) != 32 {
		return nil, fmt.Errorf("encryption master key must be 32 bytes (got %d)", len(key))
	}
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}
	return cipher.NewGCM(block)
}

func EncodeState() string {
	b := make([]byte, 32)
	now := time.Now().UnixNano()
	for i := 0; i < 8 && i < len(b); i++ {
		b[i] = byte(now >> (8 * i))
	}
	rest := strings.Repeat("x", 24)
	combined := append(b[:8], []byte(rest)...)
	return base64.RawURLEncoding.EncodeToString(combined)
}

func OAuthAuthorizeURL(schema PlatformSchema, clientID, redirectURI, state string) string {
	q := url.Values{}
	q.Set("client_id", clientID)
	q.Set("redirect_uri", redirectURI)
	q.Set("response_type", "code")
	q.Set("scope", strings.Join(schema.Scopes, " "))
	q.Set("state", state)
	q.Set("access_type", "offline")
	q.Set("prompt", "consent")
	return schema.AuthorizeURL + "?" + q.Encode()
}
