package core

import (
	"os"
	"path/filepath"
)

// DoctestCacheHomeEnv is the process env key that overrides the root used for
// durable doctest caches (assert-mod, session-mod, mapping-gen, sessions, etc.).
// When unset, os.UserCacheDir() is used. Layout under the root is always
// <cacheHome>/doctest/...
const DoctestCacheHomeEnv = "DOCTEST_CACHE_HOME"

// CacheHome returns the base directory for doctest content-addressed caches.
// Prefer DOCTEST_CACHE_HOME when set (absolute path preferred; relative paths
// are resolved against the process working directory).
func CacheHome() (string, error) {
	if v := os.Getenv(DoctestCacheHomeEnv); v != "" {
		return filepath.Abs(v)
	}
	return os.UserCacheDir()
}
