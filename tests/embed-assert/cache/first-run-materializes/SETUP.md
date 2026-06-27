# Scenario

**Feature**: first doctest test with assert import creates assert-mod cache directory

```
# cold cache (or missing content hash dir)
doctest test with assert -> creates <md5>/{assert.go,go.mod}
```

## Preconditions

- Expected cache dir for current assert source may not exist before run.

## Steps

1. Record whether expected cache dir exists.
2. Create public module with assert leaf and run `doctest test <tests> -v`.
3. Assert cache dir now exists with correct layout.

```go
import (
	"os"
	"testing"
)

var cacheExistedBefore bool

func Setup(t *testing.T, req *Request) error {
	cacheDir := expectedAssertCacheDir(t)
	if _, err := os.Stat(cacheDir); err == nil {
		cacheExistedBefore = true
		os.RemoveAll(cacheDir)
	}
	createPublicModuleProject(t, "", defaultAssertAssertGo())
	setupModuleEnv(t, req)
	req.Args = []string{"test", testDir, "-v"}
	return nil
}
```