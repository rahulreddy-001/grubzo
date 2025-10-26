package router

import (
	"grubzo/internal/config"
	"grubzo/internal/middlewares"
	"grubzo/internal/repository"
	"grubzo/internal/router/auth"
	"grubzo/internal/router/session"
	v1 "grubzo/internal/router/v1"
	"grubzo/internal/services"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Router struct {
	router *gin.Engine
	auth   *auth.Handlers
	v1     *v1.Handlers
}

func Setup(logger *zap.Logger, db *gorm.DB, repository *repository.Repository, ss *services.Services, config *config.Config) *gin.Engine {
	engine := newRouter(logger.Named("router"), db, repository, ss, config)

	engine.router.GET("/health", func(c *gin.Context) {
		c.String(200, "OK")
	})

	engine.router.NoRoute(getReactReverseProxy())

	auth := engine.router.Group("/auth")
	engine.auth.Setup(auth)

	api := engine.router.Group("/api")
	engine.v1.Setup(api)

	return engine.router
}

func getReactReverseProxy() gin.HandlerFunc {
	return func(c *gin.Context) {
		target, err := url.Parse("http://localhost:8083")
		if err != nil {
			c.String(http.StatusInternalServerError, "Invalid target URL")
			return
		}
		proxy := httputil.NewSingleHostReverseProxy(target)
		c.Request.Host = target.Host
		c.Request.URL.Host = target.Host
		c.Request.URL.Scheme = target.Scheme
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

func newRouter(logger *zap.Logger, db *gorm.DB, repository *repository.Repository, ss *services.Services, config *config.Config) *Router {
	sessionStore := session.NewMemorySessionStore()
	isDevMode := config.DevMode

	r := gin.New()
	if !isDevMode {
		gin.SetMode(gin.ReleaseMode)
	}
	r.Use(middlewares.RecoverPanic(logger.Named("painc_log")))
	r.Use(middlewares.AccessLogging(logger.Named("access_log"), isDevMode))

	authHandler := &auth.Handlers{
		Db:           db,
		Logger:       logger.Named("v1"),
		Repository:   repository,
		SessionStore: sessionStore,
		SS:           ss,
		Config:       config,
	}
	v1Handler := &v1.Handlers{
		Db:           db,
		Logger:       logger.Named("v1"),
		Repository:   repository,
		SessionStore: sessionStore,
		SS:           ss,
	}
	router := &Router{
		router: r,
		auth:   authHandler,
		v1:     v1Handler,
	}
	return router
}
