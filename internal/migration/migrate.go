package migration

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) (init bool, err error) {
	m := gormigrate.New(db, &gormigrate.Options{
		TableName:                 "migrations",
		IDColumnName:              "id",
		IDColumnSize:              190,
		UseTransaction:            false,
		ValidateUnknownMigrations: true,
	}, Migrations())
	m.InitSchema(func(db *gorm.DB) error {
		init = true
		return db.AutoMigrate(AllTables()...)
	})
	err = m.Migrate()
	return
}

func DropAll(db *gorm.DB) error {
	if err := db.Migrator().DropTable(AllTables()...); err != nil {
		return err
	}
	return db.Migrator().DropTable("migrations")
}
