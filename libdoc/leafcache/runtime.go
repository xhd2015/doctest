package leafcache

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// EnvLeafCache is the process env key for an explicit leaf-cache store root.
// When unset, the default is $CacheHome/doctest/leaf-cache/v1.
const EnvLeafCache = "DOCTEST_LEAF_CACHE"

// EnvCacheHome overrides the base cache directory (shared with other doctest caches).
const EnvCacheHome = "DOCTEST_CACHE_HOME"

// EnvSkipPaths lists warm-skip tokens (newline-separated) that the suite should
// treat as GetPass hits and skip executing. Set by the outer `doctest test`
// process after consulting the pass store.
//
// Tokens are either bare tree-relative leaf paths (single-tree suite) or
// FormatLeafIdentityEnv values (tree-qualified; multi-tree workspace / hub).
const EnvSkipPaths = "DOCTEST_LEAF_CACHE_SKIP_PATHS"

// DefaultStoreRel is the path under CacheHome for the v1 pass store.
const DefaultStoreRel = "doctest/leaf-cache/v1"

// ResolveStoreRoot returns the pass-store root:
//  1. DOCTEST_LEAF_CACHE when set
//  2. else $DOCTEST_CACHE_HOME/doctest/leaf-cache/v1 or $UserCacheDir/...
func ResolveStoreRoot() (string, error) {
	if v := strings.TrimSpace(os.Getenv(EnvLeafCache)); v != "" {
		return filepath.Abs(v)
	}
	home, err := cacheHome()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, filepath.FromSlash(DefaultStoreRel)), nil
}

func cacheHome() (string, error) {
	if v := strings.TrimSpace(os.Getenv(EnvCacheHome)); v != "" {
		return filepath.Abs(v)
	}
	return os.UserCacheDir()
}

// SkipEnabled reports whether programmatic leaf-cache skip is active for this
// Options configuration. Disable when -count is set, or -a (force), or --no-leaf-cache.
func SkipEnabled(count int, force, noLeafCache bool) bool {
	if noLeafCache || force {
		return false
	}
	// Any explicit -count=N (including 1) disables skip.
	if count > 0 {
		return false
	}
	return true
}

// KeyForLeaf builds a KeyInput for a discovered leaf under treeRoot.
// ModuleRoot is the nearest go.mod parent of treeRoot, or treeRoot if none.
func KeyForLeaf(treeRoot, leafRel, goVersion string) (KeyInput, error) {
	treeRoot, err := filepath.Abs(treeRoot)
	if err != nil {
		return KeyInput{}, err
	}
	leafDir := treeRoot
	if leafRel != "" && leafRel != "." {
		leafDir = filepath.Join(treeRoot, filepath.FromSlash(leafRel))
	}
	modRoot := findModuleRootDir(treeRoot)
	if modRoot == "" {
		modRoot = treeRoot
	}
	if goVersion == "" {
		goVersion = runtime.Version()
	}
	return KeyInput{
		ModuleRoot: modRoot,
		TreeRoot:   treeRoot,
		LeafDir:    leafDir,
		GoVersion:  goVersion,
	}, nil
}

func findModuleRootDir(dir string) string {
	dir, err := filepath.Abs(dir)
	if err != nil {
		return ""
	}
	for {
		if st, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil && !st.IsDir() {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}

// ParseSkipPaths decodes EnvSkipPaths value into a set of tree-relative paths.
func ParseSkipPaths(val string) map[string]struct{} {
	out := make(map[string]struct{})
	if val == "" {
		return out
	}
	for _, p := range strings.Split(val, "\n") {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		out[p] = struct{}{}
	}
	return out
}

// FormatSkipPaths encodes paths for EnvSkipPaths (stable sorted newlines).
func FormatSkipPaths(paths []string) string {
	if len(paths) == 0 {
		return ""
	}
	// paths are expected already unique; join as-is for deterministic order
	// callers should sort.
	return strings.Join(paths, "\n")
}
