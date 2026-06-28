package middlewares

import (
	"net"
	"net/http"
	"net/url"
	"strings"

	"grubzo/internal/utils/tenantutils"

	"github.com/gin-gonic/gin"
)

func TenantCORS(appDomain, env string, allowLocalhost bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin == "" {
			c.Next()
			return
		}

		if isAllowedOrigin(origin, c.Request.Host, appDomain, env, allowLocalhost) {
			headers := c.Writer.Header()
			headers.Set("Access-Control-Allow-Origin", origin)
			headers.Set("Access-Control-Allow-Credentials", "true")
			headers.Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			headers.Set("Access-Control-Allow-Headers", "Authorization, Content-Type, Accept, Origin, X-Requested-With")
			headers.Set("Access-Control-Max-Age", "86400")
			headers.Add("Vary", "Origin")
		}

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func isAllowedOrigin(origin, requestHost, appDomain, env string, allowLocalhost bool) bool {
	originURL, err := url.Parse(origin)
	if err != nil {
		return false
	}

	if originURL.Scheme != "http" && originURL.Scheme != "https" {
		return false
	}

	originHost := originURL.Hostname()
	if allowLocalhost && isLocalHost(originHost) {
		return true
	}

	if tenantutils.IsPlatformHost(originURL.Host, appDomain, env) {
		return true
	}

	originSubDomain, ok := tenantutils.SubDomainFromHost(originURL.Host, appDomain, env)
	if !ok {
		return false
	}

	requestSubDomain, ok := tenantutils.SubDomainFromHost(requestHost, appDomain, env)
	if !ok {
		return true
	}

	return originSubDomain == requestSubDomain
}

func isLocalHost(host string) bool {
	host = strings.ToLower(strings.TrimSpace(host))
	return host == "localhost" || net.ParseIP(host) != nil
}
