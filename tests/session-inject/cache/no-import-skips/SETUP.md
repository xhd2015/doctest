# Scenario

**Feature**: even without author session import, inject still materializes session-mod

```
# fixture tree has no github.com/xhd2015/doctest/session import in SETUP/ASSERT/Run
# but assemble always injects session.Doctest → session-mod cache is created
doctest test -> MaterializeSessionModule for inject
```

## Preconditions

- Module tree does **not** import `github.com/xhd2015/doctest/session` in author harness.
- Expected cache dir is removed before run so creation is observable.

## Steps

1. Remove expected session-mod cache dir.
2. Create public module **without** session import in fixture sources.
3. Run doctest test (product always injects `session.Doctest`).
4. Assert session-mod cache exists for the inject key.

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
