# Scenario

**Feature**: first doctest test with session import creates session-mod cache directory

```
# cold cache
doctest test with session import -> creates <md5>/{go.mod,*.go}
```

## Preconditions

- Expected cache dir for current session source may be removed before run.

## Steps

1. Remove expected session-mod cache dir if present.
2. Create public module with session-importing leaf.
3. Run `doctest test <tests> -v`.
4. Assert cache layout.

```go
import (
	"os"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	cacheDir := expectedSessionCacheDir(t)
	_ = os.RemoveAll(cacheDir)
	createPublicModuleProject(t, req, "", defaultSessionAssertGo(), true)
	setupModuleEnv(t, req)
	req.Args = []string{"test", req.TestDir, "-v"}
	return nil
}
```
