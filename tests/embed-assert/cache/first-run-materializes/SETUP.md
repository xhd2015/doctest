# Scenario

**Feature**: first doctest test with assert import creates assert-mod cache directory

```
# cold cache (isolated DOCTEST_CACHE_HOME is empty)
doctest test with assert -> creates <md5>/{assert.go,go.mod}
```

## Preconditions

- Parent `cache/SETUP.md` sets an empty isolated `DOCTEST_CACHE_HOME`.
- Expected cache dir for current assert source does not exist before run.

## Steps

1. Record that expected cache dir does not exist (cold isolated home).
2. Create public module with assert leaf and run `doctest test <tests> -v`.
3. Assert cache dir now exists with correct layout.

```go
import (
	"os"
	"testing"
)

var cacheExistedBefore bool

func Setup(t *testing.T, req *Request) error {
	cacheDir := expectedAssertCacheDir(t, req.CacheHome)
	if _, err := os.Stat(cacheDir); err == nil {
		cacheExistedBefore = true
	}
	// Do not RemoveAll: isolated DOCTEST_CACHE_HOME must never wipe the global cache.
	createPublicModuleProject(t, "", defaultAssertAssertGo())
	setupModuleEnv(t, req)
	req.Args = []string{"test", testDir, "-v"}
	return nil
}
```
