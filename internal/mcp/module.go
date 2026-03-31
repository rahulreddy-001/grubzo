package mcp

import (
	"grubzo/internal/config"
	"grubzo/internal/mcp/agent"
	"grubzo/internal/mcp/llm"
	"grubzo/internal/mcp/tools"
	internalrepo "grubzo/internal/repository"
	"grubzo/internal/router/session"
	"grubzo/internal/services"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Dependencies struct {
	Logger       *zap.Logger
	Config       *config.Config
	DB           *gorm.DB
	Redis        *redis.Client
	Repository   *internalrepo.Repository
	Services     *services.Services
	SessionStore session.Store
}

type Components struct {
	Dispatcher   *tools.Dispatcher
	AgentService *agent.Service
}

func Build(deps Dependencies) (*Components, error) {
	dispatcher := tools.NewDispatcher(
		deps.Logger.Named("tools"),
		deps.Repository,
		deps.Services,
	)
	provider, err := llm.NewProvider(deps.Config)
	if err != nil {
		return nil, err
	}
	agentService := agent.NewService(
		deps.Logger.Named("agent"),
		deps.Config,
		deps.Repository,
		dispatcher,
		provider,
	)

	return &Components{
		Dispatcher:   dispatcher,
		AgentService: agentService,
	}, nil
}
