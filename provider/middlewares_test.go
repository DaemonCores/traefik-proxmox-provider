package provider

import (
	"testing"
)

func TestParseMiddlewaresHeadersPortfolio(t *testing.T) {
	cfg := map[string]string{
		"traefik.http.middlewares.secure-headers.headers.stsseconds":            "63072000",
		"traefik.http.middlewares.secure-headers.headers.stsincludesubdomains":  "true",
		"traefik.http.middlewares.secure-headers.headers.stspreload":            "true",
		"traefik.http.middlewares.secure-headers.headers.contenttypenosniff":    "true",
		"traefik.http.middlewares.secure-headers.headers.framedeny":             "true",
		"traefik.http.middlewares.secure-headers.headers.referrerpolicy":        "strict-origin-when-cross-origin",
		"traefik.http.middlewares.secure-headers.headers.contentsecuritypolicy": "default-src 'self'",
		"traefik.http.middlewares.secure-headers.headers.customresponseheaders": "X-Robots-Tag:noindex,X-Frame-Options:DENY",
	}

	mws := parseMiddlewares(cfg)
	if len(mws) != 1 {
		t.Fatalf("expected 1 middleware, got %d", len(mws))
	}
	mw := mws["secure-headers"]
	if mw.Headers == nil {
		t.Fatal("expected a Headers middleware")
	}
	if mw.Headers.STSSeconds != 63072000 {
		t.Errorf("expected STSSeconds 63072000, got %d", mw.Headers.STSSeconds)
	}
	if !mw.Headers.STSIncludeSubdomains || !mw.Headers.STSPreload || !mw.Headers.FrameDeny || !mw.Headers.ContentTypeNosniff {
		t.Error("expected boolean header options to be true")
	}
	if mw.Headers.ReferrerPolicy != "strict-origin-when-cross-origin" {
		t.Errorf("referrerpolicy = %s", mw.Headers.ReferrerPolicy)
	}
	if mw.Headers.ContentSecurityPolicy != "default-src 'self'" {
		t.Errorf("contentsecuritypolicy = %s", mw.Headers.ContentSecurityPolicy)
	}
	if v, ok := mw.Headers.CustomResponseHeaders["X-Robots-Tag"]; !ok || v != "noindex" {
		t.Errorf("customResponseHeaders[X-Robots-Tag] = %v", v)
	}
	if v, ok := mw.Headers.CustomResponseHeaders["X-Frame-Options"]; !ok || v != "DENY" {
		t.Errorf("customResponseHeaders[X-Frame-Options] = %v", v)
	}
}

func TestParseMiddlewaresHeadersCaseInsensitive(t *testing.T) {
	// The internal label parser lowercases keys; check the camelCase form too
	mws := parseMiddlewares(map[string]string{
		"traefik.http.middlewares.h.headers.stsSeconds": "31536000",
	})
	if mws["h"] == nil || mws["h"].Headers == nil || mws["h"].Headers.STSSeconds != 31536000 {
		t.Fatalf("expected stsSeconds parsed, got %+v", mws["h"])
	}
}

func TestParseMiddlewaresRedirectRegex(t *testing.T) {
	mws := parseMiddlewares(map[string]string{
		"traefik.http.middlewares.forge-redirect.redirectregex.regex":       "^https://.*/forge$$",
		"traefik.http.middlewares.forge-redirect.redirectregex.replacement": "$${1}",
		"traefik.http.middlewares.forge-redirect.redirectregex.permanent":   "true",
	})
	mw := mws["forge-redirect"]
	if mw == nil || mw.RedirectRegex == nil {
		t.Fatal("expected a RedirectRegex middleware")
	}
	if !mw.RedirectRegex.Permanent {
		t.Error("expected permanent=true")
	}
}

func TestParseMiddlewaresStripPrefix(t *testing.T) {
	mws := parseMiddlewares(map[string]string{
		"traefik.http.middlewares.forge-stripprefix.stripprefix.prefixes": "/forge",
	})
	mw := mws["forge-stripprefix"]
	if mw == nil || mw.StripPrefix == nil {
		t.Fatal("expected a StripPrefix middleware")
	}
	if len(mw.StripPrefix.Prefixes) != 1 || mw.StripPrefix.Prefixes[0] != "/forge" {
		t.Errorf("prefixes = %v", mw.StripPrefix.Prefixes)
	}
}

func TestParseMiddlewaresMultipleTypesOnSameMiddlewareName(t *testing.T) {
	mws := parseMiddlewares(map[string]string{
		"traefik.http.middlewares.m.stripprefix.prefixes":        "/forge",
		"traefik.http.middlewares.other.redirectregex.permanent": "true",
	})
	if mws["m"] == nil || mws["m"].StripPrefix == nil {
		t.Fatal("expected stripprefix on m")
	}
	if mws["other"] == nil || mws["other"].RedirectRegex == nil {
		t.Fatal("expected redirectregex on other")
	}
}

func TestParseMiddlewaresNoneDefined(t *testing.T) {
	cfg := map[string]string{
		"traefik.enable":              "true",
		"traefik.http.routers.r.rule": "Host(`x`)",
	}
	mws := parseMiddlewares(cfg)
	if len(mws) != 0 {
		t.Fatalf("expected 0 middlewares, got %d", len(mws))
	}
}
