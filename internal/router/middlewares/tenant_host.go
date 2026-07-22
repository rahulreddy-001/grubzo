package middlewares

import (
	"net/http"

	"grubzo/internal/utils/tenantutils"

	"github.com/gin-gonic/gin"
)

func TenantHostGuard(appDomain, env string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !tenantutils.IsHostAllowedForEnv(c.Request.Host, appDomain, env) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.Next()
	}
}
