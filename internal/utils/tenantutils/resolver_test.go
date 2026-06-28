package tenantutils

import "testing"

func TestSubDomainFromHostAcceptsFEAndBEHosts(t *testing.T) {
	tests := []struct {
		name      string
		host      string
		env       string
		wantSub   string
		wantFound bool
	}{
		{name: "dev FE", host: "demo.dev.grubzo.food", env: "dev", wantSub: "demo", wantFound: true},
		{name: "dev BE", host: "demo.dev-api.grubzo.food", env: "dev", wantSub: "demo", wantFound: true},
		{name: "qa FE", host: "demo.qa.grubzo.food", env: "qa", wantSub: "demo", wantFound: true},
		{name: "qa BE", host: "demo.qa-api.grubzo.food", env: "qa", wantSub: "demo", wantFound: true},
		{name: "stage FE", host: "demo.stage.grubzo.food", env: "stage", wantSub: "demo", wantFound: true},
		{name: "stage BE", host: "demo.stage-api.grubzo.food", env: "stage", wantSub: "demo", wantFound: true},
		{name: "prod FE", host: "demo.grubzo.food", env: "prod", wantSub: "demo", wantFound: true},
		{name: "prod BE", host: "demo.api.grubzo.food", env: "prod", wantSub: "demo", wantFound: true},
		{name: "wrong env", host: "demo.qa-api.grubzo.food", env: "dev", wantFound: false},
		{name: "nested tenant", host: "a.b.dev-api.grubzo.food", env: "dev", wantFound: false},
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

func TestPlatformFromHostAcceptsFEAndBEHosts(t *testing.T) {
	tests := []struct {
		name      string
		host      string
		env       string
		wantSub   string
		wantFound bool
	}{
		{name: "dev admin FE", host: "admin.dev.grubzo.food", env: "dev", wantSub: "admin", wantFound: true},
		{name: "dev admin BE", host: "admin.dev-api.grubzo.food", env: "dev", wantSub: "admin", wantFound: true},
		{name: "prod admin FE", host: "admin.grubzo.food", env: "prod", wantSub: "admin", wantFound: true},
		{name: "prod admin BE", host: "admin.api.grubzo.food", env: "prod", wantSub: "admin", wantFound: true},
		{name: "prod platform FE", host: "platform.grubzo.food", env: "prod", wantSub: "platform", wantFound: true},
		{name: "prod platform BE", host: "platform.api.grubzo.food", env: "prod", wantSub: "platform", wantFound: true},
		{name: "tenant host", host: "demo.grubzo.food", env: "prod", wantFound: false},
		{name: "wrong env", host: "admin.qa-api.grubzo.food", env: "dev", wantFound: false},
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
		"admin.api.grubzo.food",
		"platform.dev.grubzo.food",
		"platform.dev-api.grubzo.food",
	} {
		t.Run(host, func(t *testing.T) {
			if got, ok := SubDomainFromHost(host, "grubzo.food", "dev"); ok {
				t.Fatalf("platform host returned tenant subdomain %q", got)
			}
		})
	}
}
