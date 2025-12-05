package middlewares

import (
	"grubzo/internal/repository"
	"grubzo/internal/router/ext"
	"grubzo/internal/router/session"

	"github.com/gin-gonic/gin"
)

func UserAuthenticate(repo *repository.Repository, ss session.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		sess, err := ss.GetSession(c)
		if err != nil || sess == nil || !sess.LoggedIn() {
			ss.RevokeSession(c)
			ext.Ctx(c).Unauthorized()
			return
		}
		ext.Ctx(c).SetSession(sess)
		c.Next()
	}
}
