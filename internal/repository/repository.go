package repository

import (
	"grubzo/internal/migration"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RepositoryInterface interface {
	TenantRepository
	TenantUserRepository
	TenantLocationRepository
	UserRepository
	ItemRepository
	// OrderRepository
	// OrderItemRepository
	// PaymentRepository
	FileRepository
}

var _ RepositoryInterface = (*Repository)(nil)

type Repository struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewRepository(db *gorm.DB, logger *zap.Logger, doMigration bool) (repo *Repository, init bool, err error) {
	repo = &Repository{
		db:     db,
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
