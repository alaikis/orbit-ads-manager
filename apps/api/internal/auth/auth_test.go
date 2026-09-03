package auth_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"orbit/apps/api/internal/auth"
	"orbit/apps/api/internal/model"
	"orbit/apps/api/internal/tenant"
	"orbit/apps/api/pkg/database"
	"orbit/apps/api/config"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestDB(t *testing.T) {
	config.Load()
	database.Connect()
	model.Migrate()
}

func teardownTestDB(t *testing.T) {
	_ = database.DB.Migrator().DropTable(
		&model.EmailLog{},
		&model.Notification{},
		&model.ReportSchedule{},
		&model.WebhookDelivery{},
		&model.RuleExecution{},
		&model.Rule{},
		&model.AgentAction{},
		&model.AgentToolCall{},
		&model.Message{},
		&model.Conversation{},
		&model.FeedProductSnapshot{},
		&model.FeedRun{},
		&model.Feed{},
		&model.SyncJob{},
		&model.DailyStat{},
		&model.Ad{},
		&model.AdGroup{},
		&model.Campaign{},
		&model.AdAccount{},
		&model.Order{},
		&model.ProductVariant{},
		&model.Product{},
		&model.StoreConfig{},
		&model.OAuthToken{},
		&model.PlatformConn{},
		&model.AuditLog{},
		&model.RefreshToken{},
		&model.TenantMember{},
		&model.Tenant{},
		&model.User{},
	)
}

func TestRegisterAndLogin(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.POST("/auth/register", auth.Register)
	r.POST("/auth/login", auth.Login)

	body, _ := json.Marshal(map[string]string{"email": "test@example.com", "password": "password123", "store_type": "unknown"})
	req := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	loginBody, _ := json.Marshal(map[string]string{"email": "test@example.com", "password": "password123"})
	loginReq := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(loginBody))
	loginReq.Header.Set("Content-Type", "application/json")
	loginW := httptest.NewRecorder()
	r.ServeHTTP(loginW, loginReq)
	assert.Equal(t, http.StatusOK, loginW.Code)
}

func TestHashPassword(t *testing.T) {
	hash, err := auth.HashPassword("password123")
	assert.NoError(t, err)
	assert.True(t, auth.CheckPasswordHash("password123", hash))
	assert.False(t, auth.CheckPasswordHash("wrong", hash))
}

func TestRBACMiddleware(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(auth.NewAuthMiddleware().Handle())
	r.Use(tenant.NewTenantMiddleware().Handle())
	r.GET("/test", auth.RequirePermission("store", "read"), func(c *gin.Context) { c.JSON(200, gin.H{"ok": true}) })

	token, _, _ := auth.GenerateTokenPair(1, 1, "customer_admin")
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
