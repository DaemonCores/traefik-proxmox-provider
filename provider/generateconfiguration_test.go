package provider

import (
	"testing"

	"github.com/DaemonCores/traefik-proxmox-provider/internal"
)

func TestGenerateConfigurationEmitsMiddlewares(t *testing.T) {
	labels := map[string]string{
		"traefik.enable":                                                       "true",
		"traefik.http.routers.forge.rule":                                      "Host(`forge.iartistique`) && PathPrefix(`/forge`)",
		"traefik.http.routers.forge.middlewares":                               "forge-stripprefix,forge-redirect,forge-noindex",
		"traefik.http.routers.forge.tls.certresolver":                          "le",
		"traefik.http.services.forge.loadbalancer.server.url":                  "http://forgejo.iartistique:3000",
		"traefik.http.middlewares.forge-stripprefix.stripprefix.prefixes":      "/forge",
		"traefik.http.middlewares.forge-redirect.redirectregex.regex":          "https://[^/]+/forge$$",
		"traefik.http.middlewares.forge-redirect.redirectregex.permanent":      "true",
		"traefik.http.middlewares.forge-noindex.headers.stsseconds":            "63072000",
		"traefik.http.middlewares.forge-noindex.headers.framedeny":             "true",
		"traefik.http.middlewares.forge-noindex.headers.customresponseheaders": "X-Robots-Tag:noindex",
	}

	services := map[string][]internal.Service{
		"iartistique": {internal.NewService(123, "forgejo", labels)},
	}
	config := generateConfiguration(services)

	if len(config.HTTP.Middlewares) != 3 {
		t.Fatalf("expected 3 middlewares, got %d: %v", len(config.HTTP.Middlewares), config.HTTP.Middlewares)
	}
	sp := config.HTTP.Middlewares["forge-stripprefix"].StripPrefix
	if sp == nil || len(sp.Prefixes) != 1 || sp.Prefixes[0] != "/forge" {
		t.Errorf("stripprefix = %+v", sp)
	}
	rr := config.HTTP.Middlewares["forge-redirect"].RedirectRegex
	if rr == nil || !rr.Permanent {
		t.Errorf("redirectregex = %+v", rr)
	}
	h := config.HTTP.Middlewares["forge-noindex"].Headers
	if h == nil || h.STSSeconds != 63072000 || !h.FrameDeny {
		t.Errorf("headers = %+v", h)
	}
	if v := h.CustomResponseHeaders["X-Robots-Tag"]; v != "noindex" {
		t.Errorf("customResponseHeaders = %+v", h.CustomResponseHeaders)
	}

	router := config.HTTP.Routers["forge"]
	if router == nil {
		t.Fatal("expected forge router")
	}
	if len(router.Middlewares) != 3 || router.Middlewares[0] != "forge-stripprefix" {
		t.Errorf("router middlewares = %v", router.Middlewares)
	}
}

func TestGenerateConfigurationRouterWithoutMiddlewareKeepsPayload(t *testing.T) {
	labels := map[string]string{
		"traefik.enable":                                      "true",
		"traefik.http.routers.mail.rule":                      "Host(`mail.net.gabriel-guillou.fr`)",
		"traefik.http.services.mail.loadbalancer.server.port": "3500",
	}
	services := map[string][]internal.Service{
		"iartistique": {internal.NewService(130, "synapmail", labels)},
	}
	config := generateConfiguration(services)
	if r := config.HTTP.Routers["mail"]; r == nil {
		t.Fatal("expected mail router")
	}
	if len(config.HTTP.Middlewares) != 0 {
		t.Fatalf("expected no middlewares, got %v", config.HTTP.Middlewares)
	}
}
