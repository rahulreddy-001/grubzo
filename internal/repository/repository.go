package repository

//go:generate go run ../../cmd/injecttrace -file repository.go -receiver Repository -service Repository

import (
	"grubzo/internal/migration"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RepositoryInterface interface {
	TenantRepository
	TenantUserRepository
	TenantLocationRepository
	UserRepository
	ItemRepository
	CartRepository
	WalletRepository
	OrderRepository
	FileRepository
	RoleRepository
	ChatRepository
}

var _ RepositoryInterface = (*Repository)(nil)

type Repository struct {
	db     *gorm.DB
	rdb    *redis.Client
	logger *zap.Logger
}

func NewRepository(db *gorm.DB, rdb *redis.Client, logger *zap.Logger, doMigration bool) (repo *Repository, init bool, err error) {
	repo = &Repository{
		db:     db,
		rdb:    rdb,
		logger: logger.Named("repository"),
	}
	if doMigration {
		init, err = migration.Migrate(db)
		if err != nil {
			return
		}
		if init {
			logger.Info("database schema was initialized")
		} else {
			logger.Info("database schema is up to date")
		}
	}
	return
}
