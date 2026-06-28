package repository

//go:generate go run ../../cmd/injecttrace -file utils.go -receiver Repository -service Repository
import (
	"context"

	"go.opentelemetry.io/otel"
	"gorm.io/gorm"
)

type txContextKey struct{}

func (r *Repository) dbWithContext(ctx context.Context) *gorm.DB {
	ctx, span := otel.Tracer("Repository").Start(ctx, "Repository.dbWithContext")
	defer span.End()

	if tx, ok := ctx.Value(txContextKey{}).(*gorm.DB); ok && tx != nil {
		return tx.WithContext(ctx)
	}
	return r.db.WithContext(ctx)
}

func (r *Repository) WithTransaction(ctx context.Context, fn func(context.Context) error) error {
	ctx, span := otel.Tracer("Repository").Start(ctx, "Repository.WithTransaction")
	defer span.End()

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(context.WithValue(ctx, txContextKey{}, tx))
	})
}

func (r *Repository) dbOrTx(ctx context.Context, tx *gorm.DB) *gorm.DB {
	ctx, span := otel.Tracer("Repository").Start(ctx, "Repository.dbOrTx")
	defer span.End()

	if tx != nil {
		return tx.WithContext(ctx)
	}
	return r.dbWithContext(ctx)
}
