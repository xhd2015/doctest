# Scenario

**Feature**: doctest test without session import does not create session-mod entry

```
# no session import in SETUP/ASSERT/Run
doctest test -> skip MaterializeSessionModule for session
```

## Preconditions

- Module tree does not import `github.com/xhd2015/doctest/session`.
- Expected cache dir is removed before run so absence is observable.

## Steps

1. Remove expected session-mod cache dir.
2. Create public module **without** session import.
3. Run doctest test.
4. Assert cache dir still absent.

```go
import (
	"os"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	cacheDir := expectedSessionCacheDir(t)
	_ = os.RemoveAll(cacheDir)
	createPublicModuleProject(t, "", "", false)
	setupModuleEnv(t, req)
	req.Args = []string{"test", testDir, "-v"}
	return nil
}
```
