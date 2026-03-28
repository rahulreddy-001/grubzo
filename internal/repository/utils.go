package repository

//go:generate go run ../../cmd/injecttrace -file utils.go -receiver Repository -service Repository
import (
	"context"
	"go.opentelemetry.io/otel"
	"gorm.io/gorm"
)

func (r *Repository) dbWithContext(ctx context.Context) *gorm.DB {
	ctx, span := otel.Tracer("Repository").Start(ctx, "Repository.dbWithContext")
	defer span.End()

	return r.db.WithContext(ctx)
}

func (r *Repository) dbOrTx(ctx context.Context, tx *gorm.DB) *gorm.DB {
	ctx, span := otel.Tracer("Repository").Start(ctx, "Repository.dbOrTx")
	defer span.End()

	if tx != nil {
		return tx.WithContext(ctx)
	}
	return r.dbWithContext(ctx)
}
