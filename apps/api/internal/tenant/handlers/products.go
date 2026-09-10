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

func UpdateProduct(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		httputil.BadRequest(c, "invalid id", nil)
		return
	}

	var product model.Product
	if err := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&product).Error; err != nil {
		httputil.NotFound(c, "product not found")
		return
	}

	var req struct {
		Title       *string `json:"title,omitempty"`
		Description *string `json:"description,omitempty"`
		Link        *string `json:"link,omitempty"`
		ImageURL    *string `json:"image_url,omitempty"`
		Brand       *string `json:"brand,omitempty"`
		Gtin        *string `json:"gtin,omitempty"`
		GoogleCategory *string `json:"google_product_category,omitempty"`
		ProductType *string `json:"product_type,omitempty"`
		PriceCents  *int64  `json:"price_cents,omitempty"`
		Currency    *string `json:"currency,omitempty"`
		Status      *string `json:"status,omitempty" enums:"active,inactive"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, err.Error(), nil)
		return
	}

	if req.Title != nil {
		product.Title = *req.Title
	}
	if req.Description != nil {
		product.Description = *req.Description
	}
	if req.Link != nil {
		product.Link = *req.Link
	}
	if req.ImageURL != nil {
		product.ImageURL = *req.ImageURL
	}
	if req.Brand != nil {
		product.Brand = *req.Brand
	}
	if req.Gtin != nil {
		product.Gtin = *req.Gtin
	}
	if req.GoogleCategory != nil {
		product.GoogleCategory = *req.GoogleCategory
	}
	if req.ProductType != nil {
		product.ProductType = *req.ProductType
	}
	if req.PriceCents != nil {
		product.PriceCents = *req.PriceCents
	}
	if req.Currency != nil {
		product.Currency = *req.Currency
	}
	if req.Status != nil {
		product.Status = *req.Status
	}

	if err := database.DB.Save(&product).Error; err != nil {
		httputil.InternalError(c, "failed to update product")
		return
	}
	httputil.Success(c, product)
}

func DeleteProduct(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		httputil.BadRequest(c, "invalid id", nil)
		return
	}
	if err := database.DB.Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&model.Product{}).Error; err != nil {
		httputil.InternalError(c, "failed to delete product")
		return
	}
	httputil.Success(c, gin.H{"message": "deleted"})
}

func RegisterProductRoutes(r *gin.RouterGroup) {
	products := r.Group("/products")
	products.Use(auth.NewAuthMiddleware().Handle())
	{
		products.GET("", ListProducts)
		products.GET("/:id", GetProduct)
		products.PATCH("/:id", UpdateProduct)
		products.DELETE("/:id", DeleteProduct)
	}
}
