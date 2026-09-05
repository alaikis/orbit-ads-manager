package connection

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"orbit/apps/api/config"
	"orbit/apps/api/internal/auth"
	"orbit/apps/api/internal/httputil"
	"orbit/apps/api/internal/tenant"

	"github.com/gin-gonic/gin"
)

func ListHandler(c *gin.Context) {
	tenantID, ok := tenantIDFromContext(c)
	if !ok {
		return
	}
	items, err := List(c.Request.Context(), tenantID)
	if err != nil {
		httputil.InternalError(c, "failed to list connections")
		return
	}
	httputil.Success(c, gin.H{"items": items})
}

func GetHandler(c *gin.Context) {
	tenantID, ok := tenantIDFromContext(c)
	if !ok {
		return
	}
	id, err := parseID(c.Param("id"))
	if err != nil {
		httputil.BadRequest(c, "invalid id", nil)
		return
	}
	conn, err := Get(c.Request.Context(), tenantID, id)
	if err != nil {
		httputil.NotFound(c, "connection not found")
		return
	}
	httputil.Success(c, conn)
}

func CreateHandler(c *gin.Context) {
	tenantID, ok := tenantIDFromContext(c)
	if !ok {
		return
	}
	body, _ := c.GetRawData()
	bodyStr := strings.ReplaceAll(string(body), "\\", "\"")
	var req struct {
		Platform string            `json:"platform"`
		Name     string            `json:"name"`
		Fields   map[string]string `json:"fields"`
		Scopes   []string          `json:"scopes"`
	}
	if err := json.Unmarshal([]byte(bodyStr), &req); err != nil {
		httputil.BadRequest(c, "invalid body: "+err.Error(), nil)
		return
	}
	if req.Platform == "" || req.Name == "" {
		httputil.BadRequest(c, "platform and name are required", nil)
		return
	}
	conn, err := Create(c.Request.Context(), CreateInput{
		TenantID: tenantID,
		Platform: req.Platform,
		Name:     req.Name,
		Fields:   req.Fields,
		Scopes:   req.Scopes,
	})
	if err != nil {
		if err == ErrInvalidPlatform || err == ErrAuthFlow {
			httputil.BadRequest(c, err.Error(), nil)
			return
		}
		httputil.BadRequest(c, err.Error(), nil)
		return
	}
	httputil.Created(c, conn)
}

func UpdateHandler(c *gin.Context) {
	tenantID, ok := tenantIDFromContext(c)
	if !ok {
		return
	}
	id, err := parseID(c.Param("id"))
	if err != nil {
		httputil.BadRequest(c, "invalid id", nil)
		return
	}
	body, _ := c.GetRawData()
	bodyStr := strings.ReplaceAll(string(body), "\\", "\"")
	var req struct {
		Name   *string           `json:"name"`
		Status *string           `json:"status"`
		Fields map[string]string `json:"fields"`
	}
	if err := json.Unmarshal([]byte(bodyStr), &req); err != nil {
		httputil.BadRequest(c, "invalid body: "+err.Error(), nil)
		return
	}
	conn, err := Update(c.Request.Context(), tenantID, id, UpdateInput{
		Name:   req.Name,
		Status: req.Status,
		Fields: req.Fields,
	})
	if err != nil {
		if err == ErrNotFound {
			httputil.NotFound(c, "connection not found")
			return
		}
		httputil.BadRequest(c, err.Error(), nil)
		return
	}
	httputil.Success(c, conn)
}

func DeleteHandler(c *gin.Context) {
	tenantID, ok := tenantIDFromContext(c)
	if !ok {
		return
	}
	id, err := parseID(c.Param("id"))
	if err != nil {
		httputil.BadRequest(c, "invalid id", nil)
		return
	}
	if err := Delete(c.Request.Context(), tenantID, id); err != nil {
		if err == ErrNotFound {
			httputil.NotFound(c, "connection not found")
			return
		}
		httputil.InternalError(c, "failed to delete")
		return
	}
	httputil.Success(c, gin.H{"message": "deleted"})
}

func TestHandler(c *gin.Context) {
	tenantID, ok := tenantIDFromContext(c)
	if !ok {
		return
	}
	id, err := parseID(c.Param("id"))
	if err != nil {
		httputil.BadRequest(c, "invalid id", nil)
		return
	}
	res, err := Test(c.Request.Context(), tenantID, id)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.Success(c, res)
}

func SyncHandler(c *gin.Context) {
	tenantID, ok := tenantIDFromContext(c)
	if !ok {
		return
	}
	id, err := parseID(c.Param("id"))
	if err != nil {
		httputil.BadRequest(c, "invalid id", nil)
		return
	}
	res, err := Sync(c.Request.Context(), tenantID, id)
	if err != nil {
		httputil.BadRequest(c, err.Error(), nil)
		return
	}
	httputil.Success(c, res)
}

func ListPlatformsHandler(c *gin.Context) {
	httputil.Success(c, gin.H{"items": ListPlatforms()})
}

func GoogleShoppingMerchantsHandler(c *gin.Context) {
	tenantID, ok := tenantIDFromContext(c)
	if !ok {
		return
	}
	id, err := parseID(c.Param("id"))
	if err != nil {
		httputil.BadRequest(c, "invalid id", nil)
		return
	}
	conn, err := Get(c.Request.Context(), tenantID, id)
	if err != nil {
		httputil.NotFound(c, "connection not found")
		return
	}
	if conn.Platform != "google_shopping" {
		httputil.BadRequest(c, "connection is not google_shopping", nil)
		return
	}
	tok := LoadToken(c.Request.Context(), conn.ID)
	res, err := listGoogleShoppingMerchants(c.Request.Context(), conn, tok)
	if err != nil {
		httputil.BadRequest(c, err.Error(), nil)
		return
	}
	httputil.Success(c, res)
}

func GoogleShoppingProductsHandler(c *gin.Context) {
	tenantID, ok := tenantIDFromContext(c)
	if !ok {
		return
	}
	id, err := parseID(c.Param("id"))
	if err != nil {
		httputil.BadRequest(c, "invalid id", nil)
		return
	}
	conn, err := Get(c.Request.Context(), tenantID, id)
	if err != nil {
		httputil.NotFound(c, "connection not found")
		return
	}
	if conn.Platform != "google_shopping" {
		httputil.BadRequest(c, "connection is not google_shopping", nil)
		return
	}
	merchantID := c.Query("merchant_id")
	tok := LoadToken(c.Request.Context(), conn.ID)
	res, err := listGoogleShoppingProducts(c.Request.Context(), conn, tok, merchantID)
	if err != nil {
		httputil.BadRequest(c, err.Error(), nil)
		return
	}
	httputil.Success(c, res)
}

func OAuthStartHandler(c *gin.Context) {
	tenantID, ok := tenantIDFromContext(c)
	if !ok {
		return
	}
	body, _ := c.GetRawData()
	bodyStr := strings.ReplaceAll(string(body), "\\", "\"")
	var req struct {
		Platform    string            `json:"platform"`
		Name        string            `json:"name"`
		Fields      map[string]string `json:"fields"`
		Scopes      []string          `json:"scopes"`
		RedirectTo  string            `json:"redirect_to"`
	}
	if err := json.Unmarshal([]byte(bodyStr), &req); err != nil {
		httputil.BadRequest(c, "invalid body: "+err.Error(), nil)
		return
	}
	schema, ok := GetSchema(req.Platform)
	if !ok {
		httputil.BadRequest(c, "unknown platform: "+req.Platform, nil)
		return
	}
	if schema.AuthFlow != FlowOAuth {
		httputil.BadRequest(c, "platform does not use OAuth flow", nil)
		return
	}
	prov, provErr := LoadOAuthProvider(int64(tenantID), schema.OAuthProvider)
	clientID := ""
	if provErr == nil && prov != nil {
		clientID = prov.ClientID
	} else {
		clientID = envOrConfig(schema.ClientIDEnv)
	}
	if clientID == "" {
		httputil.BadRequest(c, "OAuth client not configured (set "+schema.ClientIDEnv+")", nil)
		return
	}

	stateBytes := make([]byte, 32)
	_, _ = rand.Read(stateBytes)
	state := base64.RawURLEncoding.EncodeToString(stateBytes)

	scopes := req.Scopes
	if len(scopes) == 0 {
		scopes = schema.Scopes
	}
	os := OAuthState{
		State:       state,
		Platform:    req.Platform,
		TenantID:    tenantID,
		Name:        req.Name,
		DraftFields: nil,
		Scopes:      scopes,
		RedirectTo:  req.RedirectTo,
		ExpiresAt:   time.Now().Add(10 * time.Minute),
	}
	if err := saveOAuthState(&os); err != nil {
		httputil.InternalError(c, "failed to create state: "+err.Error())
		return
	}

	redirectBase := config.AppCfg.App.BaseURL
	if redirectBase == "" {
		redirectBase = fmt.Sprintf("%s://%s", scheme(c), c.Request.Host)
	}
	redirectURI := strings.TrimRight(redirectBase, "/") + "/api/v1/connections/oauth/callback"

	authorizeURL := OAuthAuthorizeURL(schema, clientID, redirectURI, state)
	httputil.Success(c, gin.H{
		"authorize_url": authorizeURL,
		"state":         state,
		"expires_at":    os.ExpiresAt,
	})
}

func OAuthCallbackHandler(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")
	if code == "" || state == "" {
		c.Redirect(http.StatusFound, "/settings/connections?error=missing_code_or_state")
		return
	}
	os, err := loadOAuthState(state)
	if err != nil {
		c.Redirect(http.StatusFound, "/settings/connections?error=invalid_state")
		return
	}
	if os.Used || time.Now().After(os.ExpiresAt) {
		c.Redirect(http.StatusFound, "/settings/connections?error=state_expired")
		return
	}
	schema, _ := GetSchema(os.Platform)
	prov, provErr := LoadOAuthProvider(int64(os.TenantID), schema.OAuthProvider)
	clientID := ""
	clientSecret := ""
	if provErr == nil && prov != nil {
		clientID = prov.ClientID
		clientSecret = prov.ClientSecret
	} else {
		clientID = envOrConfig(schema.ClientIDEnv)
		clientSecret = envOrConfig(schema.ClientIDEnv + "_SECRET")
	}
	redirectBase := config.AppCfg.App.BaseURL
	if redirectBase == "" {
		redirectBase = fmt.Sprintf("%s://%s", scheme(c), c.Request.Host)
	}
	redirectURI := strings.TrimRight(redirectBase, "/") + "/api/v1/connections/oauth/callback"

	accessToken, refreshToken, expiresIn, err := exchangeOAuthCode(schema, clientID, clientSecret, code, redirectURI)
	if err != nil {
		c.Redirect(http.StatusFound, "/settings/connections?error=token_exchange_failed&detail="+err.Error())
		return
	}
	conn := Connection{
		TenantID: os.TenantID,
		Type:     string(schema.AuthFlow),
		Platform: os.Platform,
		Name:     os.Name,
		Status:   "active",
		Scopes:   os.Scopes,
	}
	if err := createConnectionDirect(&conn); err != nil {
		c.Redirect(http.StatusFound, "/settings/connections?error=conn_create_failed&detail="+err.Error())
		return
	}

	// Meta: exchange short-lived token (1-2h) for long-lived (~60d) before storing.
	// Per https://developers.facebook.com/docs/facebook-login/guides/access-tokens/get-long-lived
	var expiresAt *time.Time
	if expiresIn > 0 {
		t := time.Now().Add(time.Duration(expiresIn) * time.Second)
		expiresAt = &t
	}
	if os.Platform == "meta" && accessToken != "" && clientID != "" && clientSecret != "" {
		longLived, longExpiresIn, llErr := ExchangeMetaForLongLivedToken(clientID, clientSecret, accessToken)
		if llErr != nil {
			c.Redirect(http.StatusFound, fmt.Sprintf("/settings/connections?connected=%d&warning=meta_long_lived_failed&detail=%s", conn.ID, url.QueryEscape(llErr.Error())))
			markOAuthStateUsed(state)
			return
		}
		accessToken = longLived
		if expiresAt == nil || longExpiresIn > int(time.Until(*expiresAt).Seconds()) {
			t := time.Now().Add(time.Duration(longExpiresIn) * time.Second)
			expiresAt = &t
		}
	}
	if err := StoreToken(c.Request.Context(), conn.ID, accessToken, refreshToken, expiresAt, "Bearer", os.Scopes, ""); err != nil {
		c.Redirect(http.StatusFound, "/settings/connections?error=token_store_failed&detail="+err.Error())
		return
	}
	markOAuthStateUsed(state)
	redirectTo := os.RedirectTo
	if redirectTo == "" {
		redirectTo = "/settings/connections?connected=" + fmt.Sprintf("%d", conn.ID)
	}
	c.Redirect(http.StatusFound, redirectTo)
}

func tenantIDFromContext(c *gin.Context) (uint64, bool) {
	v, exists := c.Get("tenant_id")
	if !exists {
		httputil.BadRequest(c, "tenant not resolved", nil)
		return 0, false
	}
	id, ok := v.(uint64)
	if !ok || id == 0 {
		httputil.BadRequest(c, "invalid tenant", nil)
		return 0, false
	}
	return id, true
}

func parseID(s string) (uint64, error) {
	var id uint64
	_, err := fmt.Sscanf(s, "%d", &id)
	return id, err
}

func envOrConfig(key string) string {
	if v := config.AppCfg.OAuth.GoogleClientID; key == "GOOGLE_OAUTH_CLIENT_ID" && v != "" {
		return v
	}
	return configFromOAuth(key)
}

func configFromOAuth(key string) string {
	switch key {
	case "GOOGLE_OAUTH_CLIENT_ID":
		return config.AppCfg.OAuth.GoogleClientID
	case "GOOGLE_OAUTH_CLIENT_SECRET":
		return config.AppCfg.OAuth.GoogleClientSecret
	case "META_APP_ID":
		return config.AppCfg.OAuth.MetaClientID
	case "META_APP_SECRET":
		return config.AppCfg.OAuth.MetaClientSecret
	case "BING_CLIENT_ID":
		return config.AppCfg.OAuth.BingClientID
	case "BING_CLIENT_SECRET":
		return config.AppCfg.OAuth.BingClientSecret
	}
	return ""
}

func scheme(c *gin.Context) string {
	if c.Request.TLS != nil {
		return "https"
	}
	if v := c.GetHeader("X-Forwarded-Proto"); v != "" {
		return v
	}
	return "http"
}

func RegisterConnectionRoutes(r *gin.RouterGroup) {
	conns := r.Group("/connections")
	conns.Use(auth.NewAuthMiddleware().Handle())
	conns.Use(tenant.NewTenantMiddleware().Handle())
	{
		conns.GET("", ListHandler)
		conns.POST("", CreateHandler)
		conns.GET("/platforms", ListPlatformsHandler)
		conns.GET("/:id", GetHandler)
		conns.PATCH("/:id", UpdateHandler)
		conns.DELETE("/:id", DeleteHandler)
		conns.POST("/:id/test", TestHandler)
		conns.POST("/:id/sync", SyncHandler)
		conns.POST("/oauth/start", OAuthStartHandler)
		conns.GET("/oauth/callback", OAuthCallbackHandler)
		conns.GET("/:id/google-shopping/merchants", GoogleShoppingMerchantsHandler)
		conns.GET("/:id/google-shopping/products", GoogleShoppingProductsHandler)
		conns.POST("/routes/execute", ExecuteRouteHandler)
	}
}

// ExecuteRouteRequest 路线执行请求
type ExecuteRouteRequest struct {
	SourceConnID uint64   `json:"source_conn_id" binding:"required"`
	TargetConnID uint64   `json:"target_conn_id" binding:"required"`
	Action       string   `json:"action" binding:"required"`
	Params       map[string]string `json:"params"`
}

// ExecuteRouteHandler 执行跨连接路线（如 Woo → Meta 产品同步）
// POST /api/v1/connections/routes/execute
// Body: {source_conn_id, target_conn_id, action, params}
func ExecuteRouteHandler(c *gin.Context) {
	tenantID, ok := tenantIDFromContext(c)
	if !ok {
		return
	}
	var req ExecuteRouteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "invalid body: "+err.Error(), nil)
		return
	}
	sourceConn, err := Get(c.Request.Context(), tenantID, req.SourceConnID)
	if err != nil {
		httputil.NotFound(c, "source connection not found")
		return
	}
	targetConn, err := Get(c.Request.Context(), tenantID, req.TargetConnID)
	if err != nil {
		httputil.NotFound(c, "target connection not found")
		return
	}
	sourceClient, ok := GetClient(sourceConn.Platform)
	if !ok {
		httputil.BadRequest(c, "source platform not supported: "+sourceConn.Platform, nil)
		return
	}
	targetClient, ok := GetClient(targetConn.Platform)
	if !ok {
		httputil.BadRequest(c, "target platform not supported: "+targetConn.Platform, nil)
		return
	}
	switch req.Action {
	case "sync_products":
		sourceCreds, err := loadCredentials(c.Request.Context(), sourceConn.ID)
		if err != nil {
			httputil.InternalError(c, "failed to load source credentials: "+err.Error())
			return
		}
		targetTok := loadTokenOrEmpty(c.Request.Context(), targetConn.ID)
		targetCreds, _ := loadCredentials(c.Request.Context(), targetConn.ID)
		products, ok := req.Params["products"]
		if !ok {
			syncRes, err := sourceClient.Sync(c.Request.Context(), sourceConn, nil, sourceCreds, "list_products")
			if err != nil {
				httputil.InternalError(c, "failed to list source products: "+err.Error())
				return
			}
			products, _ = syncRes["products"].(string)
		}
		var prods []map[string]string
		_ = json.Unmarshal([]byte(products), &prods)
		if len(prods) == 0 {
			httputil.BadRequest(c, "no products to sync", nil)
			return
		}
		targetCreds["products_json"] = products
		if v, ok := req.Params["catalog_id"]; ok {
			targetCreds["catalog_id"] = v
		}
		if v, ok := req.Params["merchant_id"]; ok {
			targetCreds["merchant_id"] = v
		}
		if v, ok := req.Params["target_id"]; ok {
			targetCreds["target_id"] = v
		}
		result, err := targetClient.Sync(c.Request.Context(), targetConn, targetTok, targetCreds, "upload_products")
		if err != nil {
			httputil.InternalError(c, "failed to upload products: "+err.Error())
			return
		}
		httputil.Success(c, map[string]interface{}{
			"route":    fmt.Sprintf("%s → %s", sourceConn.Platform, targetConn.Platform),
			"action":   req.Action,
			"products": len(prods),
			"result":   result,
		})
	case "create_shopping_campaign":
		targetTok := loadTokenOrEmpty(c.Request.Context(), targetConn.ID)
		targetCreds, _ := loadCredentials(c.Request.Context(), targetConn.ID)
		if targetConn.Platform != "google_ads" {
			httputil.BadRequest(c, "create_shopping_campaign requires google_ads target", nil)
			return
		}
		if v, ok := req.Params["customer_id"]; ok {
			targetCreds["customer_id"] = v
		}
		if v, ok := req.Params["campaign_name"]; ok {
			targetCreds["campaign_name"] = v
		}
		if v, ok := req.Params["daily_budget_micros"]; ok {
			targetCreds["daily_budget_micros"] = v
		}
		if v, ok := req.Params["feed_label"]; ok {
			targetCreds["feed_label"] = v
		}
		if v, ok := req.Params["merchant_id"]; ok {
			targetCreds["merchant_id"] = v
		}
		result, err := targetClient.Sync(c.Request.Context(), targetConn, targetTok, targetCreds, "create_shopping_campaign")
		if err != nil {
			httputil.InternalError(c, "failed to create shopping campaign: "+err.Error())
			return
		}
		httputil.Success(c, map[string]interface{}{
			"route":    fmt.Sprintf("%s → %s", sourceConn.Platform, targetConn.Platform),
			"action":   req.Action,
			"result":   result,
		})
	default:
		httputil.BadRequest(c, "unsupported route action: "+req.Action, nil)
	}
}
