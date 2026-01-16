package migration

import (
	"grubzo/internal/models/entity"

	"github.com/go-gormigrate/gormigrate/v2"
)

func Migrations() []*gormigrate.Migration {
	return []*gormigrate.Migration{
		v1(),
		v2(),
		v3(),
	}
}

func AllTables() []interface{} {
	return []interface{}{
		&entity.Tenant{},
		&entity.TenantLocation{},
		&entity.TenantUser{},
		&entity.User{},
		&entity.FileMeta{},
		&entity.Item{},
		&entity.UserRole{},
		&entity.WalletBalance{},
		&entity.WalletRecharge{},
		&entity.WalletTransaction{},
		&entity.Order{},
	}
}
