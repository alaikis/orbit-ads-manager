package handlers

import (
	"encoding/json"
	"fmt"
	"orbit/apps/api/internal/auth"
	"orbit/apps/api/internal/model"
	"orbit/apps/api/pkg/database"
	"orbit/apps/api/config"
	"orbit/apps/api/internal/tenant"
	"orbit/apps/api/internal/httputil"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ConnectOAuth(c *gin.Context) {
	platform := c.Param("platform")
	tenantID, _ := c.Get("tenant_id")
	userID, _ := c.Get("user_id")

	state := uuid.New().String()
	redirect := c.DefaultQuery("redirect", "/accounts")

	cfg := config.AppCfg
	var authURL string
	switch platform {
	case "google":
		authURL = fmt.Sprintf("https://accounts.google.com/o/oauth2/v2/auth?client_id=%s&redirect_uri=%s/oauth/callback&response_type=code&scope=openid+email+profile+https://www.googleapis.com/auth/adwords&access_type=offline&prompt=consent&state=%s",
			cfg.OAuth.GoogleClientID, cfg.App.APIBaseURL, state)
	case "meta":
		authURL = fmt.Sprintf("https://www.facebook.com/v19.0/dialog/oauth?client_id=%s&redirect_uri=%s/oauth/callback&scope=ads_management+ads_read+read_insights+pages_show_list&state=%s",
			cfg.OAuth.MetaClientID, cfg.App.APIBaseURL, state)
	case "bing":
		authURL = fmt.Sprintf("https://login.microsoftonline.com/common/oauth2/v2.0/authorize?client_id=%s&redirect_uri=%s/oauth/callback&response_type=code&scope=https://ads.microsoft.com/msads.manage+offline_access&state=%s",
			cfg.OAuth.BingClientID, cfg.App.APIBaseURL, state)
	default:
		httputil.BadRequest(c, "unsupported platform", nil)
		return
	}

	database.Redis.Set(c, fmt.Sprintf("oauth:state:%s", state), map[string]interface{}{
		"platform": platform,
		"tenant_id": tenantID,
		"user_id": userID,
		"redirect": redirect,
	}, 0)

	httputil.Success(c, gin.H{"authorization_url": authURL, "state": state})
}

func OAuthCallback(c *gin.Context) {
	platform := c.Query("platform")
	state := c.Query("state")
	code := c.Query("code")

	if platform == "" || state == "" || code == "" {
		httputil.BadRequest(c, "missing required parameters", nil)
		return
	}

	key := fmt.Sprintf("oauth:state:%s", state)
	stateResult, err := database.Redis.Get(c, key).Result()
	if err != nil {
		httputil.BadRequest(c, "invalid or expired state", nil)
		return
	}

	var stateData map[string]interface{}
	if err := json.Unmarshal([]byte(stateResult), &stateData); err != nil {
		httputil.BadRequest(c, "invalid state data", nil)
		return
	}
	tenantIDFloat, _ := stateData["tenant_id"].(float64)
	userIDFloat, _ := stateData["user_id"].(float64)
	tenantID := uint64(tenantIDFloat)
	adminUserID := uint64(userIDFloat)

	conn := model.PlatformConn{
		BaseModel: model.BaseModel{TenantID: tenantID},
		Type:      "ad_account",
		Platform:  platform,
		AdminUserID: adminUserID,
		Status:    "bound",
		Meta:      map[string]interface{}{"state": state, "code": code},
	}
	if err := database.DB.Create(&conn).Error; err != nil {
		httputil.InternalError(c, "failed to create connection")
		return
	}

	httputil.Success(c, gin.H{
		"message": "authorization successful",
		"platform": platform,
		"accounts": []gin.H{{"id": conn.ID, "name": fmt.Sprintf("%s Account", platform)}},
	})
}

func OAuthBind(c *gin.Context) {
	httputil.Success(c, gin.H{"message": "accounts bound successfully"})
}

func RegisterOAuthRoutes(r *gin.RouterGroup) {
	oauth := r.Group("/oauth")
	{
		oauth.GET("/connect/:platform", auth.NewAuthMiddleware().Handle(), tenant.NewTenantMiddleware().Handle(), ConnectOAuth)
		oauth.GET("/callback", OAuthCallback)
		oauth.POST("/bind", auth.NewAuthMiddleware().Handle(), tenant.NewTenantMiddleware().Handle(), OAuthBind)
	}
}
