package migration

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func v5() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "5_ADD_TENANT_SUBDOMAIN",
		Migrate: func(db *gorm.DB) error {
			if err := db.Exec("ALTER TABLE tenants ADD COLUMN IF NOT EXISTS sub_domain TEXT").Error; err != nil {
				return err
			}
			if err := db.Exec("UPDATE tenants SET sub_domain = LOWER(code) WHERE sub_domain = '' OR sub_domain IS NULL").Error; err != nil {
				return err
			}
			if err := db.Exec("ALTER TABLE tenants ALTER COLUMN sub_domain SET NOT NULL").Error; err != nil {
				return err
			}
			return db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_tenants_sub_domain ON tenants (sub_domain)").Error
		},
	}
}
