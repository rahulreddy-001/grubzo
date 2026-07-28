package router

import (
	"grubzo/internal/config"
	"grubzo/internal/mcp"
	"grubzo/internal/repository"
	"grubzo/internal/router/auth"
	mcprouter "grubzo/internal/router/mcp"
	"grubzo/internal/router/middlewares"
	"grubzo/internal/router/platform"
	"grubzo/internal/router/session"
	v1 "grubzo/internal/router/v1"
	"grubzo/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

type Router struct {
	router *gin.Engine
	auth   *auth.Handlers
	v1     *v1.Handlers
	mcp    *mcprouter.Handlers
	plat   *platform.Handlers
}

func Setup(logger *zap.Logger, db *gorm.DB, rdb *redis.Client, repository *repository.Repository, ss *services.Services, config *config.Config) *gin.Engine {
	engine := newRouter(logger.Named("router"), db, rdb, repository, ss, config)
	api := engine.router.Group("/api")

	api.GET("/health", func(c *gin.Context) {
		c.String(200, "OK")
	})
	api.GET("/metrics", gin.WrapH(promhttp.Handler()))
	api.GET("/instance", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, config.Instance)
	})

	auth := engine.router.Group("/auth")
	engine.auth.Setup(auth)

	engine.plat.Setup(engine.router)
	engine.v1.Setup(api)
	engine.mcp.Setup(engine.router)

	return engine.router
}

func newRouter(logger *zap.Logger, db *gorm.DB, rdb *redis.Client, repository *repository.Repository, ss *services.Services, config *config.Config) *Router {
	sessionStore := session.NewMemorySessionStore(config.App.Domain)
	if config.SessionStorage == "redis" {
		sessionStore = session.NewRedisSessionStore(rdb, config.App.Domain)
	}
	isDevMode := config.IsDev()

	if !isDevMode {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(middlewares.RecoverPanic(logger.Named("painc_log")))
	r.Use(middlewares.TenantCORS(config.App.Domain, config.Environment(), config.Instance, isDevMode))
	r.Use(middlewares.TenantHostGuard(config.App.Domain, config.Environment()))
	r.Use(middlewares.AccessLogging(logger.Named("access_log"), isDevMode))
	r.Use(otelgin.Middleware("grubzo_gin", otelgin.WithGinFilter(func(c *gin.Context) bool {
		return c.FullPath() != ""
	})))

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
		Config:       config,
	}

	mcpComponents, err := mcp.Build(mcp.Dependencies{
		Logger:       logger.Named("grubzo_mcp"),
		Config:       config,
		DB:           db,
		Redis:        rdb,
		Repository:   repository,
		Services:     ss,
		SessionStore: sessionStore,
	})
	if err != nil {
		logger.Fatal("failed to register mcp routes", zap.Error(err))
	}

	mcpHandlers := mcprouter.NewHandlers(
		logger.Named("grubzo_mcp_router"),
		repository,
		rdb,
		sessionStore,
		mcpComponents,
		config,
	)
	platformHandlers := platform.NewHandlers(logger.Named("platform"), repository, config)
	router := &Router{
		router: r,
		auth:   authHandler,
		v1:     v1Handler,
		mcp:    mcpHandlers,
		plat:   platformHandlers,
	}
	return router
}
