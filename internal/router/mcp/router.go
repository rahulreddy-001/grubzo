package mcp

import (
	"grubzo/internal/config"
	mcpcore "grubzo/internal/mcp"
	"grubzo/internal/mcp/actor"
	"grubzo/internal/mcp/agent"
	"grubzo/internal/mcp/tools"
	"grubzo/internal/repository"
	"grubzo/internal/router/ext"
	"grubzo/internal/router/middlewares"
	"grubzo/internal/router/session"
	"time"

	"github.com/gin-gonic/gin"
	ratelimiter "github.com/rahulreddy-001/ratelimiter/v1"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Handlers struct {
	Logger       *zap.Logger
	Repository   *repository.Repository
	SessionStore session.Store
	Tools        *tools.Dispatcher
	Agent        *agent.Service
	RDB          *redis.Client
	Config       *config.Config
}

func NewHandlers(logger *zap.Logger, repository *repository.Repository, rdb *redis.Client, sessionStore session.Store, components *mcpcore.Components, config *config.Config) *Handlers {
	return &Handlers{
		Logger:       logger,
		Repository:   repository,
		SessionStore: sessionStore,
		Tools:        components.Dispatcher,
		Agent:        components.AgentService,
		RDB:          rdb,
		Config:       config,
	}
}

func (h *Handlers) Setup(engine *gin.Engine) {
	protected := middlewares.UserAuthenticate(h.Repository, h.SessionStore, h.Config.App.Domain, h.Config.Environment(), h.Config.Instance)
	ratelimitGenerator := middlewares.RateLimiterMiddlewareGenerator()
	twoReqPerSecSlidingWindowLogForTenantAndUser := ratelimitGenerator(ratelimiter.NewSlidingWindowLog(h.RDB, 10, time.Second), middlewares.RLK_TENANT, middlewares.RLK_USER)

	api := engine.Group("/api", protected, twoReqPerSecSlidingWindowLogForTenantAndUser)
	api.POST("/chat", h.Chat)
	api.GET("/chat/sessions", h.ListChatSessions)
	api.GET("/chat/sessions/:id/messages", h.ListChatMessages)
	api.DELETE("/chat/sessions/:id", h.DeleteChatSession)

	mcp := api.Group("/mcp")
	mcp.GET("", h.Describe)
	mcp.POST("", h.HandleJSONRPC)
	mcp.GET("/tools", h.ListTools)
	mcp.POST("/tools/:tool", h.InvokeTool)
}

func actorFromGin(c *gin.Context) (actor.Actor, error) {
	userSession, err := ext.Ctx(c).GetUserSession()
	if err != nil {
		return actor.Actor{}, err
	}

	return actor.Actor{
		UserID:     userSession.UserID,
		TenantID:   userSession.TenantID,
		LocationID: userSession.Location,
		Email:      userSession.Email,
		Type:       userSession.Type,
		Roles:      append([]string(nil), userSession.Roles...),
	}, nil
}

func truncatePayloadPreview(payload any) string {
	switch typed := payload.(type) {
	case []byte:
		if len(typed) > 2048 {
			return string(typed[:2048])
		}
		return string(typed)
	case string:
		if len(typed) > 2048 {
			return typed[:2048]
		}
		return typed
	default:
		return ""
	}
}
