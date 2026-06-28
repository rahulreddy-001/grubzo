package cmd

import (
	"path/filepath"
	"testing"

	"grubzo/internal/config"
	"grubzo/internal/migration"
)

func TestDatabaseFactoryFor(t *testing.T) {
	tests := []struct {
		name   string
		dbType string
		ok     bool
	}{
		{name: "default postgres", dbType: "", ok: true},
		{name: "postgres", dbType: "postgres", ok: true},
		{name: "postgresql", dbType: "postgresql", ok: true},
		{name: "sqlite", dbType: "sqlite", ok: true},
		{name: "local", dbType: "local", ok: true},
		{name: "unknown", dbType: "oracle", ok: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := databaseFactoryFor(tt.dbType)
			if tt.ok && err != nil {
				t.Fatalf("expected factory, got error: %v", err)
			}
			if !tt.ok && err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	}
}

func TestSQLiteDatabaseFactoryOpensLocalDB(t *testing.T) {
	cfg := &config.Config{}
	cfg.Database.Type = "sqlite"
	cfg.Database.SQLite.Path = filepath.Join(t.TempDir(), "nested", "grubzo.db")

	db, err := getDatabase(cfg)
	if err != nil {
		t.Fatalf("getDatabase() error = %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("db.DB() error = %v", err)
	}
	defer sqlDB.Close()

	if _, err := migration.Migrate(db); err != nil {
		t.Fatalf("migration.Migrate() error = %v", err)
	}
}
