package connection

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
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
	clientID := envOrConfig(schema.ClientIDEnv)
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
	clientID := envOrConfig(schema.ClientIDEnv)
	clientSecret := envOrConfig(schema.ClientIDEnv + "_SECRET")
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
	var expiresAt *time.Time
	if expiresIn > 0 {
		t := time.Now().Add(time.Duration(expiresIn) * time.Second)
		expiresAt = &t
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
	}
}
