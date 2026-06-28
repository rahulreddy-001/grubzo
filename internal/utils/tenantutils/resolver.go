package tenantutils

import (
	"net"
	"strings"
)

/*
"host":   URL that the client calls (FE or BE)
"domain": base domain, e.g. "grubzo.food"
"env":    dev | qa | stage | prod

Tenant subdomain is provisioned dynamically (e.g. "demo").

FE / BE URLs per environment:

  dev:   demo.dev.grubzo.food     /  demo.dev-api.grubzo.food
  qa:    demo.qa.grubzo.food      /  demo.qa-api.grubzo.food
  stage: demo.stage.grubzo.food   /  demo.stage-api.grubzo.food
  prod:  demo.grubzo.food         /  demo.api.grubzo.food
*/

// envMiddleSegments maps each non-prod env to the FE segment that sits
// between the tenant subdomain and the base domain in the host.
// e.g. dev -> ".dev", so demo.dev.grubzo.food -> strip ".dev" -> "demo"
var envMiddleSegments = map[string]string{
	"dev":   ".dev",
	"qa":    ".qa",
	"stage": ".stage",
}

// envAPIMiddleSegments maps each non-prod env to the BE segment.
// e.g. dev -> ".dev-api", so demo.dev-api.grubzo.food -> strip ".dev-api" -> "demo"
var envAPIMiddleSegments = map[string]string{
	"dev":   ".dev-api",
	"qa":    ".qa-api",
	"stage": ".stage-api",
}

var platformSubDomains = map[string]struct{}{
	"admin":    {},
	"platform": {},
}

// SubDomainFromHost extracts the tenant subdomain from a host header.
// Returns ("", false) for invalid, bare-domain, IP, or localhost hosts.
func SubDomainFromHost(host, appDomain, env string) (string, bool) {
	sub, ok := labelFromHost(host, appDomain, env)
	if !ok {
		return "", false
	}
	if _, isPlatform := platformSubDomains[sub]; isPlatform {
		return "", false
	}
	return sub, true
}

// PlatformFromHost extracts the platform subdomain from a host header.
// Supported platform hosts follow the same FE / BE environment rules as tenants:
//
//	dev:  admin.dev.grubzo.food / admin.dev-api.grubzo.food
//	prod: admin.grubzo.food     / admin.api.grubzo.food
func PlatformFromHost(host, appDomain, env string) (string, bool) {
	sub, ok := labelFromHost(host, appDomain, env)
	if !ok {
		return "", false
	}
	if _, isPlatform := platformSubDomains[sub]; !isPlatform {
		return "", false
	}
	return sub, true
}

func IsPlatformHost(host, appDomain, env string) bool {
	_, ok := PlatformFromHost(host, appDomain, env)
	return ok
}

func labelFromHost(host, appDomain, env string) (string, bool) {
	host = strings.ToLower(strings.TrimSpace(host))
	appDomain = strings.ToLower(strings.Trim(strings.TrimSpace(appDomain), "."))
	env = strings.ToLower(strings.TrimSpace(env))

	// Strip port if present.
	if h, _, err := net.SplitHostPort(host); err == nil {
		host = h
	}

	if host == "" || appDomain == "" || host == appDomain || host == "localhost" || net.ParseIP(host) != nil {
		return "", false
	}

	// Host must end with ".<appDomain>".
	suffix := "." + appDomain
	if !strings.HasSuffix(host, suffix) {
		return "", false
	}

	// e.g. "demo.dev.grubzo.food" → "demo.dev"
	sub := strings.TrimSuffix(host, suffix)
	if sub == "" {
		return "", false
	}

	// For non-prod envs, strip either the FE segment (".dev") or the
	// matching BE segment (".dev-api").
	if seg, ok := envMiddleSegments[env]; ok {
		apiSeg := envAPIMiddleSegments[env]
		switch {
		case strings.HasSuffix(sub, apiSeg):
			sub = strings.TrimSuffix(sub, apiSeg)
		case strings.HasSuffix(sub, seg):
			sub = strings.TrimSuffix(sub, seg)
		default:
			return "", false
		}
	} else if strings.HasSuffix(sub, ".api") {
		// prod BE host: demo.api.grubzo.food -> demo
		sub = strings.TrimSuffix(sub, ".api")
	}

	// After stripping, the remaining part must be a single label (no dots).
	if sub == "" || strings.Contains(sub, ".") {
		return "", false
	}

	return sub, true
}

// HostForSubDomain builds the FE host for a given tenant subdomain and env.
//
//	dev   → demo.dev.grubzo.food
//	qa    → demo.qa.grubzo.food
//	stage → demo.stage.grubzo.food
//	prod  → demo.grubzo.food
func HostForSubDomain(subDomain, appDomain, env string) string {
	subDomain = strings.ToLower(strings.Trim(strings.TrimSpace(subDomain), "."))
	appDomain = strings.ToLower(strings.Trim(strings.TrimSpace(appDomain), "."))
	env = strings.ToLower(strings.TrimSpace(env))

	if seg, ok := envMiddleSegments[env]; ok {
		return subDomain + seg + "." + appDomain
	}
	// prod — no middle segment.
	return subDomain + "." + appDomain
}
