package middlewares

import (
	"grubzo/internal/models/query"
	"grubzo/internal/repository"
	"grubzo/internal/router/ext"
	"grubzo/internal/router/session"
	"grubzo/internal/utils/tenantutils"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
)

func UserAuthenticate(repo *repository.Repository, ss session.Store, appDomain, env string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := otel.Tracer("RouterMiddleware").Start(c.Request.Context(), "RouterMiddleware.UserAuthenticate")
		defer span.End()

		originalRequest := c.Request
		c.Request = c.Request.WithContext(ctx)
		sess, err := ss.GetSession(c)
		if err != nil || sess == nil || !sess.LoggedIn() {
			ss.RevokeSession(c)
			ext.Ctx(c).Unauthorized()
			return
		}
		if subDomain, ok := tenantutils.SubDomainFromHost(c.Request.Host, appDomain, env); ok {
			tenant, err := repo.GetTenant(ctx, query.NewTenantQuery().WithSubDomain(subDomain))
			if err != nil || tenant.ID != sess.TenantID() {
				ss.RevokeSession(c)
				ext.Ctx(c).Unauthorized()
				return
			}
		}
		c.Request = originalRequest
		ext.Ctx(c).SetSession(sess)
		c.Next()
	}
}
