package cmd

import "testing"

func TestDatabaseFactoryFor(t *testing.T) {
	tests := []struct {
		name   string
		dbType string
		ok     bool
	}{
		{name: "default postgres", dbType: "", ok: true},
		{name: "postgres", dbType: "postgres", ok: true},
		{name: "postgresql", dbType: "postgresql", ok: true},
		{name: "local file database", dbType: "local", ok: false},
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
