package repository

import (
	"context"

	"gorm.io/gorm"
)

func (r *Repository) dbWithContext(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}

func (r *Repository) dbOrTx(ctx context.Context, tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return r.dbWithContext(ctx)
}
