package auth

import (
	"encoding/json"
	"orbit/apps/api/internal/model"
	"orbit/apps/api/pkg/database"
	"orbit/apps/api/internal/httputil"
	"strings"

	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	StoreType string `json:"store_type" binding:"omitempty,oneof=woocommerce shopify unknown"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func Register(c *gin.Context) {
	var req RegisterRequest
	
	// Read raw request body
	body, err := c.GetRawData()
	if err != nil {
		httputil.BadRequest(c, "failed to read request body", nil)
		return
	}
	
	// Cloudflare Email Obfuscation workaround: replace backslash with quote
	bodyStr := strings.ReplaceAll(string(body), "\\", "\"")
	
	if err := json.Unmarshal([]byte(bodyStr), &req); err != nil {
		httputil.BadRequest(c, err.Error(), nil)
		return
	}

	var existing model.User
	if err := database.DB.Where("email = ?", req.Email).First(&existing).Error; err == nil {
		httputil.Conflict(c, "email already registered")
		return
	}

	passwordHash, err := HashPassword(req.Password)
	if err != nil {
		httputil.InternalError(c, "failed to hash password")
		return
	}

	user := model.User{
		Email:        req.Email,
		PasswordHash: passwordHash,
		Status:       "active",
	}

	if err := database.DB.Create(&user).Error; err != nil {
		httputil.InternalError(c, "failed to create user")
		return
	}

	t := model.Tenant{
		Name:   user.Email,
		Plan:   "beta",
	}
	if err := database.DB.Create(&t).Error; err != nil {
		httputil.InternalError(c, "failed to create tenant")
		return
	}

	member := model.TenantMember{
		BaseModel: model.BaseModel{TenantID: t.ID},
		UserID:    user.ID,
		Role:      "customer_admin",
		Status:    "active",
	}
	if err := database.DB.Create(&member).Error; err != nil {
		httputil.InternalError(c, "failed to create tenant member")
		return
	}

	access, refresh, err := GenerateTokenPair(user.ID, t.ID, member.Role)
	if err != nil {
		httputil.InternalError(c, "failed to generate tokens")
		return
	}

	httputil.Created(c, gin.H{
		"user": gin.H{"id": user.ID, "email": user.Email, "role": member.Role, "tenant_id": t.ID},
		"access_token": access,
		"refresh_token": refresh,
		"token_type": "Bearer",
	})
}

func Login(c *gin.Context) {
	var req LoginRequest
	
	// Read raw request body
	body, err := c.GetRawData()
	if err != nil {
		httputil.BadRequest(c, "failed to read request body", nil)
		return
	}
	
	// Cloudflare Email Obfuscation workaround: replace backslash with quote
	bodyStr := strings.ReplaceAll(string(body), "\\", "\"")
	
	if err := json.Unmarshal([]byte(bodyStr), &req); err != nil {
		httputil.BadRequest(c, err.Error(), nil)
		return
	}

	var user model.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		httputil.Unauthorized(c, "invalid credentials")
		return
	}

	if !CheckPasswordHash(req.Password, user.PasswordHash) {
		httputil.Unauthorized(c, "invalid credentials")
		return
	}

	if user.Status != "active" {
		httputil.Unauthorized(c, "account is locked or inactive")
		return
	}

	var member model.TenantMember
	if err := database.DB.Where("user_id = ? AND status = 'active'", user.ID).First(&member).Error; err != nil {
		httputil.Unauthorized(c, "no active tenant membership")
		return
	}

	access, refresh, err := GenerateTokenPair(user.ID, member.TenantID, member.Role)
	if err != nil {
		httputil.InternalError(c, "failed to generate tokens")
		return
	}

	httputil.Success(c, gin.H{
		"user": gin.H{"id": user.ID, "email": user.Email, "role": member.Role, "tenant_id": member.TenantID},
		"access_token": access,
		"refresh_token": refresh,
		"token_type": "Bearer",
	})
}

func Refresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, err.Error(), nil)
		return
	}

	claims, err := ValidateToken(req.RefreshToken)
	if err != nil {
		httputil.Unauthorized(c, "invalid refresh token")
		return
	}

	access, refresh, err := GenerateTokenPair(claims.UID, claims.TID, claims.Role)
	if err != nil {
		httputil.InternalError(c, "failed to generate tokens")
		return
	}

	httputil.Success(c, gin.H{
		"access_token": access,
		"refresh_token": refresh,
		"token_type": "Bearer",
	})
}

func Logout(c *gin.Context) {
	httputil.Success(c, gin.H{"message": "logged out"})
}

func Me(c *gin.Context) {
	userID, _ := c.Get("user_id")
	tenantID, _ := c.Get("tenant_id")
	role, _ := c.Get("role")

	user, _ := GetUserFromDB(userID.(uint64))
	if user == nil {
		httputil.NotFound(c, "user not found")
		return
	}

	httputil.Success(c, gin.H{
		"id": user.ID,
		"email": user.Email,
		"role": role,
		"tenant_id": tenantID,
		"last_login_at": user.LastLoginAt,
	})
}

func ForgotPassword(c *gin.Context) {
	httputil.Success(c, gin.H{"message": "password reset email sent if account exists"})
}

func ResetPassword(c *gin.Context) {
	httputil.Success(c, gin.H{"message": "password reset successful"})
}

func RegisterRoutes(r *gin.RouterGroup) {
	auth := r.Group("/auth")
	{
		auth.POST("/register", Register)
		auth.POST("/login", Login)
		auth.POST("/refresh", Refresh)
		auth.POST("/logout", Logout)
		auth.GET("/me", NewAuthMiddleware().Handle(), Me)
		auth.POST("/forgot-password", ForgotPassword)
		auth.POST("/reset-password", ResetPassword)
	}
}
