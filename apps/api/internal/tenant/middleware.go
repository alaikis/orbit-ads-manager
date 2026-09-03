package tenant

import (
	"orbit/apps/api/internal/auth"
	"orbit/apps/api/internal/httputil"

	"github.com/gin-gonic/gin"
)

type TenantMiddleware struct{}

func NewTenantMiddleware() *TenantMiddleware {
	return &TenantMiddleware{}
}

func (m *TenantMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		tenantID, _ := c.Get("tenant_id")

		member, err := auth.GetTenantMembership(tenantID.(uint64), userID.(uint64))
		if err != nil {
			httputil.InternalError(c, "failed to check tenant membership")
			c.Abort()
			return
		}
		if member == nil {
			httputil.Forbidden(c, "access denied to this tenant")
			c.Abort()
			return
		}

		c.Set("role", member.Role)
		c.Set("tenant_member_id", member.ID)
		c.Next()
	}
}

func (m *TenantMiddleware) IsSuperAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("role")
		if role != "super_admin" {
			httputil.Forbidden(c, "super admin only")
			c.Abort()
			return
		}
		c.Next()
	}
}
