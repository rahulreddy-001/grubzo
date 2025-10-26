package auth

import (
	"grubzo/internal/config"
	"grubzo/internal/repository"
	"grubzo/internal/router/auth/oauth"
	"grubzo/internal/router/auth/oauth/github"
	"grubzo/internal/router/auth/oauth/google"
	"grubzo/internal/router/session"
	"grubzo/internal/services"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Handlers struct {
	Db           *gorm.DB
	Logger       *zap.Logger
	Repository   *repository.Repository
	SessionStore session.Store
	SS           *services.Services
	Config       *config.Config
}

func (h Handlers) Setup(r *gin.RouterGroup) {
	api := r.Group("/v1")
	{
		api.GET("/me", h.Me)
		api.POST("/login", h.Login)
		api.POST("/logout", h.Logout)

		oauthGroup := api.Group("/oauth")
		// path: `/auth/v1/oauth/login/<provider>`, callback: `/auth/v1/oauth/login/<provider>/callback`
		oauth.New().SetProviders(
			google.Init(
				h.Config.OAuthCreds[`google`].ClientId,
				h.Config.OAuthCreds[`google`].ClientSecret,
				h.Config.OAuthCreds[`google`].CallBackURL,
			),
			github.Init(
				h.Config.OAuthCreds[`github`].ClientId,
				h.Config.OAuthCreds[`github`].ClientSecret,
				h.Config.OAuthCreds[`github`].CallBackURL,
			),
		).UseRouter(oauthGroup).WithSessionStore(h.SessionStore).WithRepository(h.Repository).WithLogger(h.Logger).Init()
	}
}
