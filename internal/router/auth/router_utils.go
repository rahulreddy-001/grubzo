package auth

import (
	"grubzo/internal/router/ext"
	"grubzo/internal/services/rbac/permission"
	"grubzo/internal/services/rbac/role"
	"slices"

	"github.com/gin-gonic/gin"
)

func (h Handlers) isAdmin(c *gin.Context) bool {
	userSession, _ := ext.Ctx(c).GetUserSession()
	return slices.Contains(userSession.Roles, role.Admin)
}

func (h Handlers) isUser(c *gin.Context) bool {
	userSession, _ := ext.Ctx(c).GetUserSession()
	return userSession.Type == role.User
}

func (h Handlers) hasAccessTo(c *gin.Context, perm permission.Permission) bool {
	tenantID := ext.Ctx(c).TenantID()
	userSession, _ := ext.Ctx(c).GetUserSession()
	return h.SS.RBAC.IsAnyGranted(tenantID, userSession.Roles, perm)
}
