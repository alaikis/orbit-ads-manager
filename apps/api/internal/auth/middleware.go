package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct{}

func NewAuthMiddleware() *AuthMiddleware {
	return &AuthMiddleware{}
}

func (m *AuthMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 1001, "message": "missing authorization header"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 1001, "message": "invalid authorization header format"})
			c.Abort()
			return
		}

		claims, err := ValidateToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 1002, "message": "invalid or expired token"})
			c.Abort()
			return
		}

		user, err := GetUserFromDB(claims.UID)
		if err != nil || user == nil || user.Status != "active" {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 1001, "message": "user not found or inactive"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UID)
		c.Set("tenant_id", claims.TID)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func (m *AuthMiddleware) Optional() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}
		m.Handle()(c)
	}
}

func RequirePermission(resource, action string) gin.HandlerFunc {
	roleMatrix := map[string]map[string]map[string]bool{
		"super_admin": {
			"tenant":       {"read": true, "update": true},
			"member":       {"read": true, "invite": true, "remove": true, "change_role": true},
			"store":        {"read": true, "create": true, "update": true, "delete": true, "connect": true},
			"ad_account":   {"read": true, "connect": true, "update": true, "delete": true},
			"ad_campaign":  {"read": true, "update": true},
			"feed":         {"read": true, "create": true, "update": true, "regenerate": true, "delete": true},
			"rule":         {"read": true, "create": true, "update": true, "enable": true, "disable": true, "delete": true},
			"rule_execution":{"read": true},
			"report":       {"read": true, "export": true, "schedule": true},
			"agent":        {"chat": true, "approve_action": true, "reject_action": true},
			"notification": {"read": true, "mark_read": true},
			"audit_log":    {"read": true},
		},
		"customer_admin": {
			"tenant":       {"read": true, "update": true},
			"member":       {"read": true, "invite": true, "remove": true, "change_role": true},
			"store":        {"read": true, "create": true, "update": true, "delete": true, "connect": true},
			"ad_account":   {"read": true, "connect": true, "update": true, "delete": true},
			"ad_campaign":  {"read": true, "update": true},
			"feed":         {"read": true, "create": true, "update": true, "regenerate": true, "delete": true},
			"rule":         {"read": true, "create": true, "update": true, "enable": true, "disable": true, "delete": true},
			"rule_execution":{"read": true},
			"report":       {"read": true, "export": true, "schedule": true},
			"agent":        {"chat": true, "approve_action": true, "reject_action": true},
			"notification": {"read": true, "mark_read": true},
		},
		"operator": {
			"member":       {"read": true},
			"store":        {"read": true, "create": true, "update": true},
			"ad_account":   {"read": true, "update": true},
			"ad_campaign":  {"read": true, "update": true},
			"feed":         {"read": true},
			"rule":         {"read": true, "create": true, "update": true, "enable": true, "disable": true},
			"rule_execution":{"read": true},
			"report":       {"read": true, "export": true},
			"agent":        {"chat": true},
			"notification": {"read": true, "mark_read": true},
		},
		"viewer": {
			"member":       {"read": true},
			"store":        {"read": true},
			"ad_account":   {"read": true},
			"ad_campaign":  {"read": true},
			"feed":         {"read": true},
			"rule":         {"read": true},
			"rule_execution":{"read": true},
			"report":       {"read": true},
			"agent":        {"chat": true},
			"notification": {"read": true, "mark_read": true},
		},
	}

	return func(c *gin.Context) {
		role, _ := c.Get("role")
		roleStr, _ := role.(string)
		permissions, ok := roleMatrix[roleStr]
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"code": 1002, "message": "forbidden"})
			c.Abort()
			return
		}
		resPerms, ok := permissions[resource]
		if !ok || !resPerms[action] {
			c.JSON(http.StatusForbidden, gin.H{"code": 1002, "message": "forbidden"})
			c.Abort()
			return
		}
		c.Next()
	}
}
