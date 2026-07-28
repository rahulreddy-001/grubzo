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

Tenant URLs per environment:

  dev:   demo.dev.grubzo.food
  qa:    demo.qa.grubzo.food
  stage: demo.stage.grubzo.food
  prod:  demo.grubzo.food

Backend traffic is routed on the same host by path prefixes such as
/api, /platform, and /auth.
*/

// envMiddleSegments maps each non-prod env to the FE segment that sits
// between the tenant subdomain and the base domain in the host.
// e.g. dev -> ".dev", so demo.dev.grubzo.food -> strip ".dev" -> "demo"
var envMiddleSegments = map[string]string{
	"dev":   ".dev",
	"qa":    ".qa",
	"stage": ".stage",
}

func normalizeEnvironment(env string) (string, bool) {
	switch strings.ToLower(strings.TrimSpace(env)) {
	case "", "dev", "development":
		return "dev", true
	case "qa":
		return "qa", true
	case "stage", "staging":
		return "stage", true
	case "prod", "production":
		return "prod", true
	default:
		return "", false
	}
}

func normalizeInstance(instance string) string {
	return strings.ToLower(strings.Trim(strings.TrimSpace(instance), "."))
}

// SubDomainFromHost extracts the tenant subdomain from a host header.
// Returns ("", false) for invalid, bare-domain, IP, or localhost hosts.
func SubDomainFromHost(host, appDomain, env, instance string) (string, bool) {
	sub, ok := labelFromHost(host, appDomain, env)
	if !ok {
		return "", false
	}
	if instance = normalizeInstance(instance); instance != "" && sub == instance {
		return "", false
	}
	return sub, true
}

// PlatformFromHost extracts the platform subdomain from a host header.
// Supported platform hosts follow the same environment rules as tenants:
//
//	dev:  <instance>.dev.grubzo.food
//	prod: <instance>.grubzo.food
func PlatformFromHost(host, appDomain, env, instance string) (string, bool) {
	sub, ok := labelFromHost(host, appDomain, env)
	if !ok {
		return "", false
	}
	if instance = normalizeInstance(instance); instance == "" || sub != instance {
		return "", false
	}
	return sub, true
}

func IsPlatformHost(host, appDomain, env, instance string) bool {
	_, ok := PlatformFromHost(host, appDomain, env, instance)
	return ok
}

// IsHostAllowedForEnv returns true when a request host is either outside the
// managed app domain (for localhost/IP based development) or matches the
// configured environment exactly.
func IsHostAllowedForEnv(host, appDomain, env string) bool {
	host = strings.ToLower(strings.TrimSpace(host))
	appDomain = strings.ToLower(strings.Trim(strings.TrimSpace(appDomain), "."))

	if h, _, err := net.SplitHostPort(host); err == nil {
		host = h
	}

	if host == "" || appDomain == "" || host == "localhost" || net.ParseIP(host) != nil {
		return true
	}

	if host != appDomain && !strings.HasSuffix(host, "."+appDomain) {
		return true
	}

	_, ok := labelFromHost(host, appDomain, env)
	return ok
}

func labelFromHost(host, appDomain, env string) (string, bool) {
	host = strings.ToLower(strings.TrimSpace(host))
	appDomain = strings.ToLower(strings.Trim(strings.TrimSpace(appDomain), "."))
	env, ok := normalizeEnvironment(env)
	if !ok {
		return "", false
	}

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

	// For non-prod envs, strip the environment segment (".dev", ".qa", ".stage").
	if seg, ok := envMiddleSegments[env]; ok {
		if strings.HasSuffix(sub, seg) {
			sub = strings.TrimSuffix(sub, seg)
		} else {
			return "", false
		}
	} else if env != "prod" {
		return "", false
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
