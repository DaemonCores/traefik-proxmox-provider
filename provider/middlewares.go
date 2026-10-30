package provider

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/DaemonCores/traefik-proxmox-provider/internal"
	"github.com/traefik/genconf/dynamic"
)

// parseMiddlewares builds dynamic middleware definitions from labels like
// traefik.http.middlewares.<name>.<type>.<field>=value.
// Labels must already be normalised to lowercase (the internal parser does
// this), except for custom header maps whose keys are preserved verbatim.
func parseMiddlewares(cfg map[string]string) map[string]*dynamic.Middleware {
	middlewares := make(map[string]*dynamic.Middleware)

	// Defensive lowercase normalisation: the internal parser already
	// lowercases keys, keep parseMiddlewares correct on raw input too.
	normalized := make(map[string]string, len(cfg))
	for k, v := range cfg {
		normalized[strings.ToLower(k)] = v
	}
	cfg = normalized

	names := make(map[string]bool)
	for k := range cfg {
		if strings.HasPrefix(k, "traefik.http.middlewares.") {
			parts := strings.Split(k, ".")
			if len(parts) > 3 {
				names[parts[3]] = true
			}
		}
	}

	for _, name := range mapKeysToSlice(names) {
		prefix := fmt.Sprintf("traefik.http.middlewares.%s", name)

		if h := parseHeadersMiddleware(cfg, prefix); h != nil {
			middlewares[name] = &dynamic.Middleware{Headers: h}
		}
		if r := parseRedirectRegexMiddleware(cfg, prefix); r != nil {
			middlewares[name] = &dynamic.Middleware{RedirectRegex: r}
		}
		if s := parseStripPrefixMiddleware(cfg, prefix); s != nil {
			middlewares[name] = &dynamic.Middleware{StripPrefix: s}
		}
	}

	return middlewares
}

func parseHeadersMiddleware(cfg map[string]string, prefix string) *dynamic.Headers {
	labels, hasType := anyValuePrefix(cfg, prefix+".headers")
	if !hasType {
		return nil
	}

	headers := &dynamic.Headers{}

	// String fields
	stringFields := map[string]*string{
		"sslhost":                 &headers.SSLHost,
		"customframeoptionsvalue": &headers.CustomFrameOptionsValue,
		"custombrowserxssvalue":   &headers.CustomBrowserXSSValue,
		"contentsecuritypolicy":   &headers.ContentSecurityPolicy,
		"publickey":               &headers.PublicKey,
		"referrerpolicy":          &headers.ReferrerPolicy,
		"featurepolicy":           &headers.FeaturePolicy,
		"permissionspolicy":       &headers.PermissionsPolicy,
	}
	for label, dst := range stringFields {
		if val, exists := labels["."+label]; exists && val != "" {
			*dst = val
		}
	}

	// Boolean fields
	boolFields := map[string]*bool{
		"sslredirect":                   &headers.SSLRedirect,
		"ssltemporaryredirect":          &headers.SSLTemporaryRedirect,
		"sslforcehost":                  &headers.SSLForceHost,
		"stsincludesubdomains":          &headers.STSIncludeSubdomains,
		"stspreload":                    &headers.STSPreload,
		"forcestsheader":                &headers.ForceSTSHeader,
		"framedeny":                     &headers.FrameDeny,
		"contenttypenosniff":            &headers.ContentTypeNosniff,
		"browserxssfilter":              &headers.BrowserXSSFilter,
		"accesscontrolallowcredentials": &headers.AccessControlAllowCredentials,
		"addvaryheader":                 &headers.AddVaryHeader,
	}
	for label, dst := range boolFields {
		if val, exists := labels["."+label]; exists {
			if b, err := stringToBool(val); err == nil {
				*dst = b
			}
		}
	}

	// Integer fields
	if val, exists := labels[".stsseconds"]; exists && val != "" {
		if n, err := strconv.ParseInt(val, 10, 64); err == nil {
			headers.STSSeconds = n
		}
	}
	if val, exists := labels[".accesscontrolmaxage"]; exists && val != "" {
		if n, err := strconv.ParseInt(val, 10, 64); err == nil {
			headers.AccessControlMaxAge = n
		}
	}

	// List fields
	listFields := map[string]*[]string{
		"allowedhosts":                      &headers.AllowedHosts,
		"hostsproxyheaders":                 &headers.HostsProxyHeaders,
		"accesscontrolallowheaders":         &headers.AccessControlAllowHeaders,
		"accesscontrolallowmethods":         &headers.AccessControlAllowMethods,
		"accesscontrolalloworiginlist":      &headers.AccessControlAllowOriginList,
		"accesscontrolalloworiginlistregex": &headers.AccessControlAllowOriginListRegex,
		"accesscontrolexposeheaders":        &headers.AccessControlExposeHeaders,
	}
	for label, dst := range listFields {
		if val, exists := labels["."+label]; exists && val != "" {
			*dst = splitList(val)
		}
	}

	// Map fields (key1:value1,key2:value2)
	if val, exists := labels[".customrequestheaders"]; exists && val != "" {
		headers.CustomRequestHeaders = parseKeyValueList(val)
	}
	if val, exists := labels[".customresponseheaders"]; exists && val != "" {
		headers.CustomResponseHeaders = parseKeyValueList(val)
	}
	if val, exists := labels[".sslproxyheaders"]; exists && val != "" {
		headers.SSLProxyHeaders = parseKeyValueList(val)
	}

	if isEmptyHeaders(headers) {
		return nil
	}
	return headers
}

func isEmptyHeaders(h *dynamic.Headers) bool {
	return len(h.CustomRequestHeaders) == 0 &&
		len(h.CustomResponseHeaders) == 0 &&
		len(h.SSLProxyHeaders) == 0 &&
		h.AccessControlAllowCredentials == false &&
		len(h.AccessControlAllowHeaders) == 0 &&
		len(h.AccessControlAllowMethods) == 0 &&
		len(h.AccessControlAllowOriginList) == 0 &&
		len(h.AccessControlAllowOriginListRegex) == 0 &&
		len(h.AccessControlExposeHeaders) == 0 &&
		h.AccessControlMaxAge == 0 &&
		h.AddVaryHeader == false &&
		len(h.AllowedHosts) == 0 &&
		len(h.HostsProxyHeaders) == 0 &&
		h.SSLRedirect == false &&
		h.SSLTemporaryRedirect == false &&
		h.SSLHost == "" &&
		h.SSLForceHost == false &&
		h.STSSeconds == 0 &&
		h.STSIncludeSubdomains == false &&
		h.STSPreload == false &&
		h.ForceSTSHeader == false &&
		h.FrameDeny == false &&
		h.CustomFrameOptionsValue == "" &&
		h.ContentTypeNosniff == false &&
		h.BrowserXSSFilter == false &&
		h.CustomBrowserXSSValue == "" &&
		h.ContentSecurityPolicy == "" &&
		h.PublicKey == "" &&
		h.ReferrerPolicy == "" &&
		h.FeaturePolicy == "" &&
		h.PermissionsPolicy == "" &&
		h.IsDevelopment == false
}

func parseRedirectRegexMiddleware(cfg map[string]string, prefix string) *dynamic.RedirectRegex {
	labels, hasType := anyValuePrefix(cfg, prefix+".redirectregex")
	if !hasType {
		return nil
	}
	rr := &dynamic.RedirectRegex{}
	if val, exists := labels[".regex"]; exists {
		rr.Regex = val
	}
	if val, exists := labels[".replacement"]; exists {
		rr.Replacement = val
	}
	if val, exists := labels[".permanent"]; exists {
		if b, err := stringToBool(val); err == nil {
			rr.Permanent = b
		}
	}
	return rr
}

func parseStripPrefixMiddleware(cfg map[string]string, prefix string) *dynamic.StripPrefix {
	labels, hasType := anyValuePrefix(cfg, prefix+".stripprefix")
	if !hasType {
		return nil
	}
	sp := &dynamic.StripPrefix{}
	if val, exists := labels[".prefixes"]; exists && val != "" {
		sp.Prefixes = splitList(val)
	}
	if val, exists := labels[".forceslash"]; exists {
		if b, err := stringToBool(val); err == nil {
			sp.ForceSlash = b
		}
	}
	return sp
}

// anyValuePrefix returns the sub-labels of a prefix (keys with the prefix
// removed) and reports whether any label with that prefix exists at all.
func anyValuePrefix(cfg map[string]string, prefix string) (map[string]string, bool) {
	sub := make(map[string]string)
	found := false
	for k, v := range cfg {
		if strings.HasPrefix(k, prefix+".") {
			sub[strings.TrimPrefix(k, prefix)] = v
			found = true
		}
	}
	return sub, found
}

func splitList(val string) []string {
	parts := strings.Split(val, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		out = append(out, strings.TrimSpace(p))
	}
	return out
}

func parseKeyValueList(val string) map[string]string {
	out := make(map[string]string)
	for _, pair := range strings.Split(val, ",") {
		parts := strings.SplitN(pair, ":", 2)
		key := strings.TrimSpace(parts[0])
		if key == "" {
			continue
		}
		value := key
		if len(parts) == 2 {
			value = strings.TrimSpace(parts[1])
		}
		out[key] = value
	}
	return out
}

// merge middleware definitions extracted per-service into the global config,
// sorted by name for deterministic output. First definition wins on a name
// clash so repeated CT annotations never rewrite an earlier one.
func mergeMiddlewares(into map[string]*dynamic.Middleware, service internal.Service) {
	extracted := parseMiddlewares(service.Config)
	names := make([]string, 0, len(extracted))
	for name := range extracted {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		if _, exists := into[name]; !exists {
			into[name] = extracted[name]
		}
	}
}
