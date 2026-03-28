package middlewares

import (
	"net/http"

	"grubzo/internal/router/ext"
	"grubzo/internal/router/session"
	"grubzo/internal/services/rbac"
	"grubzo/internal/services/rbac/permission"
	"grubzo/internal/utils"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
)

func AccessControlMiddlewareGenerator(r *rbac.RBAC, ss session.Store) func(perms ...permission.Permission) gin.HandlerFunc {
	return func(perms ...permission.Permission) gin.HandlerFunc {
		return func(c *gin.Context) {
			_, span := otel.Tracer("RouterMiddleware").Start(c.Request.Context(), "RouterMiddleware.AccessControl")
			defer span.End()

			info, err := ext.Ctx(c).GetUserSession()
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "Invalid session",
				})
				return
			}
			userSession, err := utils.AsType[session.UserSession](info)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "Invalid session format",
				})
				return
			}

			for _, p := range perms {
				if !r.IsAnyGranted(userSession.TenantID, userSession.Roles, p) {
					c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
						"error": "Access Denied",
					})
					return
				}
			}
			c.Next()
		}
	}
}
