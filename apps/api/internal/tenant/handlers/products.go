package handlers

import (
	"strconv"

	"orbit/apps/api/internal/auth"
	"orbit/apps/api/internal/model"
	"orbit/apps/api/pkg/database"
	"orbit/apps/api/internal/httputil"

	"github.com/gin-gonic/gin"
)

func ListProducts(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")
	search := c.Query("search")

	query := database.DB.Where("tenant_id = ? AND deleted_at IS NULL", tenantID)
	if search != "" {
		query = query.Where("title ILIKE ?", "%"+search+"%")
	}

	var products []model.Product
	if err := query.Preload("Variants").Limit(100).Find(&products).Error; err != nil {
		httputil.InternalError(c, "failed to list products")
		return
	}
	httputil.Success(c, gin.H{"items": products})
}

func GetProduct(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		httputil.BadRequest(c, "invalid id", nil)
		return
	}

	var product model.Product
	if err := database.DB.Preload("Variants").Where("id = ? AND tenant_id = ?", id, tenantID).First(&product).Error; err != nil {
		httputil.NotFound(c, "product not found")
		return
	}
	httputil.Success(c, product)
}

func RegisterProductRoutes(r *gin.RouterGroup) {
	products := r.Group("/products")
	products.Use(auth.NewAuthMiddleware().Handle())
	{
		products.GET("", ListProducts)
		products.GET("/:id", GetProduct)
	}
}
