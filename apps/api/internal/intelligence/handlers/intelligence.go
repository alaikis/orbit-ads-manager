package handlers

import (
	"fmt"

	"orbit/apps/api/internal/auth"
	"orbit/apps/api/internal/intelligence"
	"orbit/apps/api/internal/model"
	"orbit/apps/api/pkg/database"
	"orbit/apps/api/internal/tenant"
	"orbit/apps/api/internal/httputil"

	"github.com/gin-gonic/gin"
)

func AnalyzeAudience(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")
	analyzer := intelligence.NewAudienceAnalyzer()
	segments, err := analyzer.Analyze(c.Request.Context(), tenantID.(uint64))
	if err != nil {
		httputil.InternalError(c, fmt.Sprintf("failed to analyze audience: %v", err))
		return
	}
	httputil.Success(c, gin.H{"segments": segments})
}

func GenerateStrategy(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")
	productID := c.Param("productId")
	if productID == "" {
		httputil.BadRequest(c, "product_id is required", nil)
		return
	}

	var productIDUint uint64
	if _, err := fmt.Sscanf(productID, "%d", &productIDUint); err != nil {
		httputil.BadRequest(c, "invalid product_id", nil)
		return
	}

	engine := intelligence.NewStrategyEngine(nil)
	strategy, err := engine.Generate(c.Request.Context(), tenantID.(uint64), productIDUint)
	if err != nil {
		httputil.InternalError(c, fmt.Sprintf("failed to generate strategy: %v", err))
		return
	}
	httputil.Success(c, strategy)
}

func GenerateMultiPlatformStrategy(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")
	engine := intelligence.NewStrategyEngine(nil)
	strategies, err := engine.GenerateMultiPlatform(c.Request.Context(), tenantID.(uint64))
	if err != nil {
		httputil.InternalError(c, fmt.Sprintf("failed to generate strategies: %v", err))
		return
	}
	httputil.Success(c, gin.H{"strategies": strategies})
}

func OptimizeBudget(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")
	var req struct {
		TotalBudget float64 `json:"total_budget" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, err.Error(), nil)
		return
	}

	engine := intelligence.NewStrategyEngine(nil)
	budgetSplit, err := engine.OptimizeBudget(c.Request.Context(), tenantID.(uint64), req.TotalBudget)
	if err != nil {
		httputil.InternalError(c, fmt.Sprintf("failed to optimize budget: %v", err))
		return
	}
	httputil.Success(c, budgetSplit)
}

func GenerateCreative(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")
	var req struct {
		ProductID uint64 `json:"product_id" binding:"required"`
		Platform  string `json:"platform" binding:"required"`
		Type      string `json:"type" binding:"required,oneof=text image video"`
		Style     string `json:"style"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, err.Error(), nil)
		return
	}

	var product model.Product
	if err := database.DB.Where("id = ? AND tenant_id = ?", req.ProductID, tenantID).First(&product).Error; err != nil {
		httputil.NotFound(c, "product not found")
		return
	}

	generator := intelligence.NewCreativeGenerator(nil)
	genReq := intelligence.GenerationRequest{
		Type:        req.Type,
		ProductName: product.Title,
		Description: product.Description,
		Category:    product.ProductType,
		Platform:    req.Platform,
		Language:    "en",
		Variations:  3,
		Style:       req.Style,
	}

	var result interface{}
	var err error

	switch req.Type {
	case "text":
		creatives, err := generator.GenerateText(c.Request.Context(), genReq)
		if err != nil {
			httputil.InternalError(c, fmt.Sprintf("failed to generate creative: %v", err))
			return
		}
		result = map[string]interface{}{"creatives": creatives}
	case "image":
		prompt := fmt.Sprintf("Professional product photo of %s, clean background, e-commerce style", product.Title)
		imageURL, err := generator.GenerateImage(c.Request.Context(), prompt)
		if err != nil {
			httputil.InternalError(c, fmt.Sprintf("failed to generate image: %v", err))
			return
		}
		result = map[string]interface{}{"image_url": imageURL, "prompt": prompt}
	case "video":
		prompt := fmt.Sprintf("Product showcase video of %s, 15 seconds, professional lighting", product.Title)
		videoURL, err := generator.GenerateVideo(c.Request.Context(), prompt, 15)
		if err != nil {
			httputil.InternalError(c, fmt.Sprintf("failed to generate video: %v", err))
			return
		}
		result = map[string]interface{}{"video_url": videoURL, "prompt": prompt}
	default:
		httputil.BadRequest(c, "unsupported creative type", nil)
		return
	}

	if err != nil {
		httputil.InternalError(c, fmt.Sprintf("failed to generate creative: %v", err))
		return
	}

	httputil.Success(c, result)
}

func GetRecommendations(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")
	engine := intelligence.NewRecommendationEngine(nil)
	recommendations, err := engine.GenerateRecommendations(c.Request.Context(), tenantID.(uint64))
	if err != nil {
		httputil.InternalError(c, fmt.Sprintf("failed to generate recommendations: %v", err))
		return
	}
	httputil.Success(c, gin.H{"recommendations": recommendations})
}

func GetBidRecommendation(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")
	campaignID := c.Param("campaignId")
	if campaignID == "" {
		httputil.BadRequest(c, "campaign_id is required", nil)
		return
	}

	var campaignIDUint uint64
	if _, err := fmt.Sscanf(campaignID, "%d", &campaignIDUint); err != nil {
		httputil.BadRequest(c, "invalid campaign_id", nil)
		return
	}

	optimizer := intelligence.NewBiddingOptimizer(nil)
	bid, err := optimizer.Optimize(c.Request.Context(), tenantID.(uint64), campaignIDUint)
	if err != nil {
		httputil.InternalError(c, fmt.Sprintf("failed to optimize bidding: %v", err))
		return
	}
	httputil.Success(c, bid)
}

func RegisterIntelligenceRoutes(r *gin.RouterGroup) {
	intel := r.Group("/intelligence")
	intel.Use(auth.NewAuthMiddleware().Handle())
	intel.Use(tenant.NewTenantMiddleware().Handle())
	{
		intel.GET("/audience", AnalyzeAudience)
		intel.GET("/strategy/:productId", GenerateStrategy)
		intel.GET("/strategy/multi", GenerateMultiPlatformStrategy)
		intel.POST("/budget/optimize", OptimizeBudget)
		intel.POST("/creative", GenerateCreative)
		intel.GET("/recommendations", GetRecommendations)
		intel.GET("/bidding/:campaignId", GetBidRecommendation)
	}
}
