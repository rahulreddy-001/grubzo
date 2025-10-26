package migration

import (
	"grubzo/internal/models/entity"

	"github.com/go-gormigrate/gormigrate/v2"
)

func Migrations() []*gormigrate.Migration {
	return []*gormigrate.Migration{
		v1(),
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
	}
}
