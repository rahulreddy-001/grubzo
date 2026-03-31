package mcp

import (
	mcpcore "grubzo/internal/mcp"
	"grubzo/internal/mcp/actor"
	"grubzo/internal/mcp/agent"
	"grubzo/internal/mcp/tools"
	"grubzo/internal/repository"
	"grubzo/internal/router/ext"
	"grubzo/internal/router/middlewares"
	"grubzo/internal/router/session"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handlers struct {
	Logger       *zap.Logger
	Repository   *repository.Repository
	SessionStore session.Store
	Tools        *tools.Dispatcher
	Agent        *agent.Service
}

func NewHandlers(logger *zap.Logger, repository *repository.Repository, sessionStore session.Store, components *mcpcore.Components) *Handlers {
	return &Handlers{
		Logger:       logger,
		Repository:   repository,
		SessionStore: sessionStore,
		Tools:        components.Dispatcher,
		Agent:        components.AgentService,
	}
}

func (h *Handlers) Setup(engine *gin.Engine) {
	protected := middlewares.UserAuthenticate(h.Repository, h.SessionStore)

	api := engine.Group("/api", protected)
	api.POST("/chat", h.Chat)
	api.GET("/chat/sessions", h.ListChatSessions)
	api.GET("/chat/sessions/:id/messages", h.ListChatMessages)
	api.DELETE("/chat/sessions/:id", h.DeleteChatSession)

	mcp := engine.Group("/mcp", protected)
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
