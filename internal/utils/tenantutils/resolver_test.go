package tenantutils

import "testing"

func TestSubDomainFromHostAcceptsEnvironmentHosts(t *testing.T) {
	tests := []struct {
		name      string
		host      string
		env       string
		wantSub   string
		wantFound bool
	}{
		{name: "dev", host: "demo.dev.grubzo.food", env: "dev", wantSub: "demo", wantFound: true},
		{name: "qa", host: "demo.qa.grubzo.food", env: "qa", wantSub: "demo", wantFound: true},
		{name: "stage", host: "demo.stage.grubzo.food", env: "stage", wantSub: "demo", wantFound: true},
		{name: "prod", host: "demo.grubzo.food", env: "prod", wantSub: "demo", wantFound: true},
		{name: "old dev api host", host: "demo.dev-api.grubzo.food", env: "dev", wantFound: false},
		{name: "old prod api host", host: "demo.api.grubzo.food", env: "prod", wantFound: false},
		{name: "wrong env", host: "demo.qa.grubzo.food", env: "dev", wantFound: false},
		{name: "nested tenant", host: "a.b.dev.grubzo.food", env: "dev", wantFound: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSub, gotFound := SubDomainFromHost(tt.host, "grubzo.food", tt.env)
			if gotFound != tt.wantFound {
				t.Fatalf("found = %v, want %v", gotFound, tt.wantFound)
			}
			if gotSub != tt.wantSub {
				t.Fatalf("subdomain = %q, want %q", gotSub, tt.wantSub)
			}
		})
	}
}

func TestSubDomainFromHostRejectsWrongEnvironmentCombinations(t *testing.T) {
	tests := []struct {
		name string
		host string
		env  string
	}{
		{name: "dev rejects prod host", host: "demo.grubzo.food", env: "dev"},
		{name: "dev rejects qa host", host: "demo.qa.grubzo.food", env: "dev"},
		{name: "dev rejects stage host", host: "demo.stage.grubzo.food", env: "dev"},
		{name: "qa rejects prod host", host: "demo.grubzo.food", env: "qa"},
		{name: "qa rejects dev host", host: "demo.dev.grubzo.food", env: "qa"},
		{name: "qa rejects stage host", host: "demo.stage.grubzo.food", env: "qa"},
		{name: "stage rejects prod host", host: "demo.grubzo.food", env: "stage"},
		{name: "stage rejects dev host", host: "demo.dev.grubzo.food", env: "stage"},
		{name: "stage rejects qa host", host: "demo.qa.grubzo.food", env: "stage"},
		{name: "prod rejects dev host", host: "demo.dev.grubzo.food", env: "prod"},
		{name: "prod rejects qa host", host: "demo.qa.grubzo.food", env: "prod"},
		{name: "prod rejects stage host", host: "demo.stage.grubzo.food", env: "prod"},
		{name: "unknown env rejects host", host: "demo.grubzo.food", env: "demo.dev.grubzo.food"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, ok := SubDomainFromHost(tt.host, "grubzo.food", tt.env); ok {
				t.Fatalf("subdomain = %q, want not found", got)
			}
		})
	}
}

func TestIsHostAllowedForEnv(t *testing.T) {
	tests := []struct {
		name string
		host string
		env  string
		want bool
	}{
		{name: "dev tenant", host: "demo.dev.grubzo.food", env: "dev", want: true},
		{name: "development alias", host: "demo.dev.grubzo.food", env: "development", want: true},
		{name: "qa tenant", host: "demo.qa.grubzo.food", env: "qa", want: true},
		{name: "stage tenant", host: "demo.stage.grubzo.food", env: "stage", want: true},
		{name: "staging alias", host: "demo.stage.grubzo.food", env: "staging", want: true},
		{name: "prod tenant", host: "demo.grubzo.food", env: "prod", want: true},
		{name: "production alias", host: "demo.grubzo.food", env: "production", want: true},
		{name: "dev rejects prod host", host: "demo.grubzo.food", env: "dev", want: false},
		{name: "prod rejects dev host", host: "demo.dev.grubzo.food", env: "prod", want: false},
		{name: "dev rejects api host", host: "demo.dev-api.grubzo.food", env: "dev", want: false},
		{name: "prod rejects api host", host: "demo.api.grubzo.food", env: "prod", want: false},
		{name: "bare domain rejected", host: "grubzo.food", env: "prod", want: false},
		{name: "localhost allowed", host: "localhost:8082", env: "dev", want: true},
		{name: "ip allowed", host: "127.0.0.1:8082", env: "dev", want: true},
		{name: "outside domain allowed", host: "internal.service", env: "dev", want: true},
		{name: "unknown env rejected", host: "demo.grubzo.food", env: "demo.dev.grubzo.food", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsHostAllowedForEnv(tt.host, "grubzo.food", tt.env); got != tt.want {
				t.Fatalf("allowed = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlatformFromHostAcceptsEnvironmentHosts(t *testing.T) {
	tests := []struct {
		name      string
		host      string
		env       string
		wantSub   string
		wantFound bool
	}{
		{name: "dev admin", host: "admin.dev.grubzo.food", env: "dev", wantSub: "admin", wantFound: true},
		{name: "prod admin", host: "admin.grubzo.food", env: "prod", wantSub: "admin", wantFound: true},
		{name: "prod platform", host: "platform.grubzo.food", env: "prod", wantSub: "platform", wantFound: true},
		{name: "tenant host", host: "demo.grubzo.food", env: "prod", wantFound: false},
		{name: "old dev api host", host: "admin.dev-api.grubzo.food", env: "dev", wantFound: false},
		{name: "old prod api host", host: "admin.api.grubzo.food", env: "prod", wantFound: false},
		{name: "wrong env", host: "admin.qa.grubzo.food", env: "dev", wantFound: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSub, gotFound := PlatformFromHost(tt.host, "grubzo.food", tt.env)
			if gotFound != tt.wantFound {
				t.Fatalf("found = %v, want %v", gotFound, tt.wantFound)
			}
			if gotSub != tt.wantSub {
				t.Fatalf("subdomain = %q, want %q", gotSub, tt.wantSub)
			}
		})
	}
}

func TestSubDomainFromHostRejectsPlatformHosts(t *testing.T) {
	for _, host := range []string{
		"admin.grubzo.food",
		"platform.dev.grubzo.food",
	} {
		t.Run(host, func(t *testing.T) {
			if got, ok := SubDomainFromHost(host, "grubzo.food", "dev"); ok {
				t.Fatalf("platform host returned tenant subdomain %q", got)
			}
		})
	}
}
