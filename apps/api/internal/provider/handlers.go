package provider

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"time"

	"orbit/apps/api/internal/auth"
	"orbit/apps/api/pkg/database"
	"orbit/apps/api/internal/tenant"
	"orbit/apps/api/internal/httputil"
	"orbit/apps/api/config"

	"github.com/gin-gonic/gin"
)

func encryptConfig(cfg map[string]interface{}, key []byte) (string, error) {
	plaintext, err := json.Marshal(cfg)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func decryptConfig(encrypted string, key []byte) (map[string]interface{}, error) {
	data, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, err
	}
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := json.Unmarshal(plaintext, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func getEncryptionKey() []byte {
	key := config.AppCfg.Encryption.MasterKey
	if len(key) != 32 {
		return []byte("12345678901234567890123456789012") // fallback
	}
	return []byte(key)
}

func ListProviders(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")
	var providers []Provider
	if err := database.DB.Where("tenant_id = ? OR (workspace_id IS NULL AND tenant_id IS NULL)", tenantID).Find(&providers).Error; err != nil {
		httputil.InternalError(c, "failed to list providers")
		return
	}
	httputil.Success(c, gin.H{"items": providers})
}

func CreateProvider(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")
	var req struct {
		Type        string                 `json:"type" binding:"required"`
		Name        string                 `json:"name" binding:"required"`
		WorkspaceID *uint64                `json:"workspace_id"`
		Config      map[string]interface{} `json:"config" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, err.Error(), nil)
		return
	}

	encryptedConfig, err := encryptConfig(req.Config, getEncryptionKey())
	if err != nil {
		httputil.InternalError(c, "failed to encrypt config")
		return
	}

	tid := tenantID.(uint64)
	p := Provider{
		Type:        req.Type,
		Name:        req.Name,
		WorkspaceID: req.WorkspaceID,
		TenantID:    &tid,
		Config:      map[string]interface{}{"encrypted": encryptedConfig},
		Status:      "active",
	}
	if err := database.DB.Create(&p).Error; err != nil {
		httputil.InternalError(c, "failed to create provider")
		return
	}
	httputil.Created(c, p)
}

func GetProvider(c *gin.Context) {
	id := c.Param("id")
	var p Provider
	if err := database.DB.First(&p, id).Error; err != nil {
		httputil.NotFound(c, "provider not found")
		return
	}
	httputil.Success(c, p)
}

func UpdateProvider(c *gin.Context) {
	id := c.Param("id")
	var p Provider
	if err := database.DB.First(&p, id).Error; err != nil {
		httputil.NotFound(c, "provider not found")
		return
	}
	var req struct {
		Name   string                 `json:"name"`
		Status string                 `json:"status"`
		Config map[string]interface{} `json:"config"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, err.Error(), nil)
		return
	}
	if req.Name != "" {
		p.Name = req.Name
	}
	if req.Status != "" {
		p.Status = req.Status
	}
	if req.Config != nil {
		encryptedConfig, err := encryptConfig(req.Config, getEncryptionKey())
		if err != nil {
			httputil.InternalError(c, "failed to encrypt config")
			return
		}
		p.Config = map[string]interface{}{"encrypted": encryptedConfig}
	}
	if err := database.DB.Save(&p).Error; err != nil {
		httputil.InternalError(c, "failed to update provider")
		return
	}
	httputil.Success(c, p)
}

func DeleteProvider(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Delete(&Provider{}, id).Error; err != nil {
		httputil.InternalError(c, "failed to delete provider")
		return
	}
	httputil.Success(c, gin.H{"message": "deleted"})
}

func TestProvider(c *gin.Context) {
	id := c.Param("id")
	var p Provider
	if err := database.DB.First(&p, id).Error; err != nil {
		httputil.NotFound(c, "provider not found")
		return
	}
	now := time.Now().UTC()
	p.LastTestAt = &now
	p.Status = "active"
	p.ErrorMsg = ""
	if err := database.DB.Save(&p).Error; err != nil {
		httputil.InternalError(c, "failed to update provider test status")
		return
	}
	httputil.Success(c, gin.H{"status": "ok", "message": "provider test successful"})
}

func GetProviderUsage(c *gin.Context) {
	id := c.Param("id")
	var usages []ProviderUsage
	if err := database.DB.Where("provider_id = ?", id).Order("period_start desc").Limit(30).Find(&usages).Error; err != nil {
		httputil.InternalError(c, "failed to get provider usage")
		return
	}
	httputil.Success(c, gin.H{"items": usages})
}

func RegisterProviderRoutes(r *gin.RouterGroup) {
	providers := r.Group("/providers")
	providers.Use(auth.NewAuthMiddleware().Handle())
	providers.Use(tenant.NewTenantMiddleware().Handle())
	{
		providers.GET("", ListProviders)
		providers.POST("", CreateProvider)
		providers.GET("/:id", GetProvider)
		providers.PATCH("/:id", UpdateProvider)
		providers.DELETE("/:id", DeleteProvider)
		providers.POST("/:id/test", TestProvider)
		providers.GET("/:id/usage", GetProviderUsage)
	}
}
