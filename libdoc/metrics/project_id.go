package metrics

import (
	"crypto/sha256"
	"encoding/hex"
	"net/url"
	"strings"
)

// ProjectIDFromOrigin normalizes a git remote origin URL into a stable slug.
// Scheme, userinfo, and trailing ".git" are stripped; path separators become "_".
//
// Examples:
//
//	https://github.com/xhd2015/doctest.git → github.com_xhd2015_doctest
//	git@github.com:xhd2015/doctest.git     → github.com_xhd2015_doctest
func ProjectIDFromOrigin(origin string) string {
	origin = strings.TrimSpace(origin)
	if origin == "" {
		return ""
	}

	var host, path string

	// scp-like: git@host:path or user@host:path (no scheme)
	if !strings.Contains(origin, "://") {
		// ssh:// is handled below; bare scp form uses host:path
		if at := strings.Index(origin, "@"); at >= 0 {
			rest := origin[at+1:]
			if colon := strings.Index(rest, ":"); colon >= 0 && !strings.Contains(rest[:colon], "/") {
				host = rest[:colon]
				path = rest[colon+1:]
				return slugify(host, path)
			}
		}
		// host:path without user
		if colon := strings.Index(origin, ":"); colon >= 0 && !strings.Contains(origin[:colon], "/") {
			host = origin[:colon]
			path = origin[colon+1:]
			return slugify(host, path)
		}
	}

	// URL forms: https://..., ssh://..., git://...
	u, err := url.Parse(origin)
	if err == nil && u.Host != "" {
		host = u.Host
		// strip optional port
		if h, _, ok := strings.Cut(host, ":"); ok && !strings.Contains(host, "]") {
			// keep IPv6 in brackets; for simple host:port drop port
			if !strings.Contains(h, ":") {
				host = h
			}
		}
		path = strings.TrimPrefix(u.Path, "/")
		return slugify(host, path)
	}

	// Fallback: strip scheme-like prefix and treat remainder as host/path
	s := origin
	if i := strings.Index(s, "://"); i >= 0 {
		s = s[i+3:]
	}
	if at := strings.Index(s, "@"); at >= 0 {
		s = s[at+1:]
	}
	s = strings.TrimPrefix(s, "/")
	s = strings.TrimSuffix(s, ".git")
	s = strings.ReplaceAll(s, "/", "_")
	s = strings.ReplaceAll(s, ":", "_")
	return s
}

func slugify(host, path string) string {
	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")
	if strings.HasSuffix(strings.ToLower(path), ".git") {
		path = path[:len(path)-len(".git")]
	}
	// drop userinfo leftovers in host
	if at := strings.LastIndex(host, "@"); at >= 0 {
		host = host[at+1:]
	}
	combined := host
	if path != "" {
		combined = host + "/" + path
	}
	combined = strings.ReplaceAll(combined, "/", "_")
	return combined
}

// ProjectIDFallback returns nogit_<first 12 hex chars of sha256(absRoot)>.
func ProjectIDFallback(absRoot string) string {
	sum := sha256.Sum256([]byte(absRoot))
	return "nogit_" + hex.EncodeToString(sum[:])[:12]
}
