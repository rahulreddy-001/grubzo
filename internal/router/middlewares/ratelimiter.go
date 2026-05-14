package middlewares

import (
	"grubzo/internal/router/ext"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	ratelimiter "github.com/rahulreddy-001/ratelimiter/v1"
	"go.opentelemetry.io/otel"
)

type RateLimiterKeys int

const (
	RLK_IP     RateLimiterKeys = 1
	RLK_URL    RateLimiterKeys = 2
	RLK_TENANT RateLimiterKeys = 3
	RLK_USER   RateLimiterKeys = 4
)

func RateLimiterMiddlewareGenerator() func(strategy ratelimiter.Ratelimiter, keys ...RateLimiterKeys) gin.HandlerFunc {
	return func(strategy ratelimiter.Ratelimiter, keys ...RateLimiterKeys) gin.HandlerFunc {
		return func(c *gin.Context) {
			ctx, span := otel.Tracer("RouterMiddleware").Start(c.Request.Context(), "RouterMiddleware.RateLimiterMiddlewareGenerator")
			defer span.End()

			parts := []string{}

			for _, keyType := range keys {
				switch keyType {
				case RLK_IP:
					parts = append(parts, c.ClientIP())

				case RLK_URL:
					parts = append(parts, c.Request.URL.Path)

				case RLK_TENANT:
					parts = append(parts, strconv.Itoa(int(ext.Ctx(c).TenantID())))

				case RLK_USER:
					parts = append(parts, strconv.Itoa(int(ext.Ctx(c).UserID())))
				}
			}

			key := strings.Join(parts, ":")

			if !strategy.Consume(ctx, key, 1) {
				c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
					"error": "Rate limit exceeded. Please try again later.",
				})
				return
			}
		}
	}
}
