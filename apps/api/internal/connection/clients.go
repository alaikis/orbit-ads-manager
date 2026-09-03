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

func listBingAdAccounts(ctx context.Context, conn *Connection, creds map[string]string) (map[string]interface{}, error) {
	clientID := creds["client_id"]
	clientSecret := creds["client_secret"]
	devToken := creds["developer_token"]
	if clientID == "" || clientSecret == "" || devToken == "" {
		return nil, fmt.Errorf("missing bing credentials (client_id, client_secret, developer_token)")
	}
	tok := creds["access_token"]
	if tok == "" {
		return nil, fmt.Errorf("bing sync requires OAuth authorization; use /connections/oauth/start")
	}
	_ = clientID
	_ = clientSecret
	_ = devToken
	return map[string]interface{}{"ok": true, "message": "bing ad account listing requires v1.1 (System User setup); credentials stored"}, nil
}

func listTikTokAdAccounts(ctx context.Context, conn *Connection, creds map[string]string) (map[string]interface{}, error) {
	tok := creds["access_token"]
	if tok == "" {
		return nil, fmt.Errorf("tiktok sync requires OAuth authorization; use /connections/oauth/start")
	}
	return map[string]interface{}{"ok": true, "message": "tiktok advertiser listing requires v1.1 (full OAuth pipeline); credentials stored"}, nil
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
