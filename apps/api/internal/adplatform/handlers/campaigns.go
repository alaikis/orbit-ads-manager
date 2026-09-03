package handlers

import (
	"fmt"
	"orbit/apps/api/internal/adplatform"
	"orbit/apps/api/internal/auth"
	"orbit/apps/api/internal/crypto"
	"orbit/apps/api/internal/model"
	"orbit/apps/api/pkg/database"
	"orbit/apps/api/internal/tenant"
	"orbit/apps/api/internal/httputil"
	"orbit/apps/api/config"

	"github.com/gin-gonic/gin"
)

func getAdapterForAccount(c *gin.Context, accountID uint64) (adplatform.PlatformAdapter, adplatform.PlatformConfig, error) {
	var account model.AdAccount
	if err := database.DB.Preload("Conn").Where("id = ? AND tenant_id = ?", accountID, c.MustGet("tenant_id")).First(&account).Error; err != nil {
		return nil, adplatform.PlatformConfig{}, fmt.Errorf("ad account not found: %w", err)
	}

	var conn model.PlatformConn
	if err := database.DB.Where("id = ?", account.ConnID).First(&conn).Error; err != nil {
		return nil, adplatform.PlatformConfig{}, fmt.Errorf("platform connection not found: %w", err)
	}

	var oauthToken model.OAuthToken
	if err := database.DB.Where("conn_id = ?", conn.ID).First(&oauthToken).Error; err != nil {
		return nil, adplatform.PlatformConfig{}, fmt.Errorf("oauth token not found: %w", err)
	}

	accessToken, err := crypto.Decrypt(oauthToken.EncryptedToken1)
	if err != nil {
		return nil, adplatform.PlatformConfig{}, fmt.Errorf("failed to decrypt access token: %w", err)
	}

	cfg := config.AppCfg
	var adapter adplatform.PlatformAdapter
	switch conn.Platform {
	case "meta":
		adapter = adplatform.NewMetaAdapter(cfg)
	case "google":
		adapter = adplatform.NewGoogleAdapter(cfg)
	default:
		return nil, adplatform.PlatformConfig{}, fmt.Errorf("unsupported platform: %s", conn.Platform)
	}

	platformConfig := adplatform.PlatformConfig{
		AccessToken: accessToken,
		AccountID:   account.ExternalID,
		CustomerID:  account.ExternalID,
	}

	return adapter, platformConfig, nil
}

func ListCampaigns(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")
	var accounts []model.AdAccount
	if err := database.DB.Where("tenant_id = ? AND deleted_at IS NULL", tenantID).Find(&accounts).Error; err != nil {
		httputil.InternalError(c, "failed to list ad accounts")
		return
	}

	var allCampaigns []gin.H
	for _, account := range accounts {
		adapter, platformConfig, err := getAdapterForAccount(c, account.ID)
		if err != nil {
			continue
		}

		campaigns, err := adapter.ListCampaigns(c, platformConfig)
		if err != nil {
			continue
		}

		for _, camp := range campaigns {
			allCampaigns = append(allCampaigns, gin.H{
				"id":            camp.ID,
				"name":          camp.Name,
				"status":        camp.Status,
				"daily_budget":  camp.DailyBudget,
				"platform":      account.Platform,
				"account_id":    account.ID,
				"external_id":   camp.ID,
			})
		}
	}

	httputil.Success(c, gin.H{"items": allCampaigns})
}

func GetCampaign(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")
	id := c.Param("id")

	var campaign model.Campaign
	if err := database.DB.Preload("AdAccount").Where("id = ? AND tenant_id = ?", id, tenantID).First(&campaign).Error; err != nil {
		httputil.NotFound(c, "campaign not found")
		return
	}

	adapter, platformConfig, err := getAdapterForAccount(c, campaign.AdAccountID)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}

	fetched, err := adapter.GetCampaign(c, platformConfig, campaign.ExternalID)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}

	campaign.Status = fetched.Status
	campaign.Name = fetched.Name
	campaign.DailyBudgetCents = int64(fetched.DailyBudget * 100)
	database.DB.Save(&campaign)

	httputil.Success(c, campaign)
}

func PauseCampaign(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")
	id := c.Param("id")

	var campaign model.Campaign
	if err := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&campaign).Error; err != nil {
		httputil.NotFound(c, "campaign not found")
		return
	}

	adapter, platformConfig, err := getAdapterForAccount(c, campaign.AdAccountID)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}

	if err := adapter.UpdateCampaignStatus(c, platformConfig, campaign.ExternalID, "paused"); err != nil {
		httputil.InternalError(c, err.Error())
		return
	}

	campaign.Status = "paused"
	database.DB.Save(&campaign)
	httputil.Success(c, gin.H{"message": "campaign paused", "campaign_id": campaign.ID})
}

func ResumeCampaign(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")
	id := c.Param("id")

	var campaign model.Campaign
	if err := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&campaign).Error; err != nil {
		httputil.NotFound(c, "campaign not found")
		return
	}

	adapter, platformConfig, err := getAdapterForAccount(c, campaign.AdAccountID)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}

	if err := adapter.UpdateCampaignStatus(c, platformConfig, campaign.ExternalID, "enabled"); err != nil {
		httputil.InternalError(c, err.Error())
		return
	}

	campaign.Status = "enabled"
	database.DB.Save(&campaign)
	httputil.Success(c, gin.H{"message": "campaign resumed", "campaign_id": campaign.ID})
}

func UpdateBudget(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")
	id := c.Param("id")
	var req struct {
		NewBudgetCents int64 `json:"new_budget_cents" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, err.Error(), nil)
		return
	}

	var campaign model.Campaign
	if err := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&campaign).Error; err != nil {
		httputil.NotFound(c, "campaign not found")
		return
	}

	adapter, platformConfig, err := getAdapterForAccount(c, campaign.AdAccountID)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}

	if err := adapter.UpdateCampaignBudget(c, platformConfig, campaign.ExternalID, req.NewBudgetCents); err != nil {
		httputil.InternalError(c, err.Error())
		return
	}

	campaign.DailyBudgetCents = req.NewBudgetCents
	database.DB.Save(&campaign)
	httputil.Success(c, gin.H{"message": "budget updated", "campaign_id": campaign.ID, "new_budget_cents": req.NewBudgetCents})
}

func ListAdGroups(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")
	campaignID := c.Param("id")

	var campaign model.Campaign
	if err := database.DB.Where("id = ? AND tenant_id = ?", campaignID, tenantID).First(&campaign).Error; err != nil {
		httputil.NotFound(c, "campaign not found")
		return
	}

	adapter, platformConfig, err := getAdapterForAccount(c, campaign.AdAccountID)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}

	groups, err := adapter.ListAdGroups(c, platformConfig, campaign.ExternalID)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}

	result := make([]gin.H, 0, len(groups))
	for _, g := range groups {
		result = append(result, gin.H{
			"id":          g.ID,
			"name":        g.Name,
			"status":      g.Status,
			"campaign_id": campaign.ID,
		})
	}

	httputil.Success(c, gin.H{"items": result})
}

func ListAds(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")
	adGroupID := c.Param("id")

	var adGroup model.AdGroup
	if err := database.DB.Where("id = ? AND tenant_id = ?", adGroupID, tenantID).First(&adGroup).Error; err != nil {
		httputil.NotFound(c, "ad group not found")
		return
	}

	var campaign model.Campaign
	if err := database.DB.Where("id = ? AND tenant_id = ?", adGroup.CampaignID, tenantID).First(&campaign).Error; err != nil {
		httputil.NotFound(c, "campaign not found")
		return
	}

	adapter, platformConfig, err := getAdapterForAccount(c, campaign.AdAccountID)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}

	ads, err := adapter.ListAds(c, platformConfig, adGroup.ExternalID)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}

	result := make([]gin.H, 0, len(ads))
	for _, ad := range ads {
		result = append(result, gin.H{
			"id":          ad.ID,
			"name":        ad.Name,
			"status":      ad.Status,
			"ad_group_id": adGroup.ID,
		})
	}

	httputil.Success(c, gin.H{"items": result})
}

func RegisterCampaignRoutes(r *gin.RouterGroup) {
	campaigns := r.Group("/advertising/campaigns")
	campaigns.Use(auth.NewAuthMiddleware().Handle())
	campaigns.Use(tenant.NewTenantMiddleware().Handle())
	{
		campaigns.GET("", auth.RequirePermission("ad_campaign", "read"), ListCampaigns)
		campaigns.GET("/:id", auth.RequirePermission("ad_campaign", "read"), GetCampaign)
		campaigns.POST("/:id/pause", auth.RequirePermission("ad_campaign", "update"), PauseCampaign)
		campaigns.POST("/:id/resume", auth.RequirePermission("ad_campaign", "update"), ResumeCampaign)
		campaigns.POST("/:id/budget", auth.RequirePermission("ad_campaign", "update"), UpdateBudget)
		campaigns.GET("/:id/ad-groups", auth.RequirePermission("ad_campaign", "read"), ListAdGroups)
	}

	ads := r.Group("/advertising/ads")
	ads.Use(auth.NewAuthMiddleware().Handle())
	ads.Use(tenant.NewTenantMiddleware().Handle())
	{
		ads.GET("", auth.RequirePermission("ad_campaign", "read"), func(c *gin.Context) {
			adGroupID := c.Query("ad_group_id")
			tenantID, _ := c.Get("tenant_id")
			var results []model.Ad
			db := database.DB.Preload("AdGroup")
			if adGroupID != "" {
				db = db.Where("ad_group_id = ? AND tenant_id = ?", adGroupID, tenantID)
			}
			if err := db.Find(&results).Error; err != nil {
				httputil.InternalError(c, "failed to list ads")
				return
			}
			httputil.Success(c, gin.H{"items": results})
		})
	}

	adGroups := r.Group("/advertising/ad-groups")
	adGroups.Use(auth.NewAuthMiddleware().Handle())
	adGroups.Use(tenant.NewTenantMiddleware().Handle())
	{
		adGroups.GET("", auth.RequirePermission("ad_campaign", "read"), ListAdGroups)
		adGroups.GET("/:id/ads", auth.RequirePermission("ad_campaign", "read"), ListAds)
		adGroups.POST("/:id/bid", auth.RequirePermission("ad_campaign", "update"), func(c *gin.Context) {
			httputil.Success(c, gin.H{"message": "bid updated", "ad_group_id": c.Param("id")})
		})
	}
}
