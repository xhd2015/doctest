# Scenario

**Feature**: `doctest test --cold-cache` on a tiny fixture exits 0 with isolated GOCACHE announce

```
# cold-cache smoke
temp module + DOCTEST_CACHE_HOME sandbox
doctest test --cold-cache mytest/
  -> exit 0
  -> stderr mentions cold-cache and GOCACHE (opts.GoCache isolation)
```

## Preconditions

- Isolated `DOCTEST_CACHE_HOME` on subprocess env (never parent Setenv).
- Tiny pass leaf (no session Assert required for this smoke).
- Deep gen-dir modes remain in `tests/test/cold-cache/`.

## Steps

1. Create temp module fixture.
2. Sandbox cache home on `req.Env`.
3. Run `doctest test --cold-cache <testDir>`.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	_, testDir := createTempModuleFixture(t, false)
	req.FixtureDir = testDir
	req.CacheHome = t.TempDir()
	req.Env = append(req.Env, "DOCTEST_CACHE_HOME="+req.CacheHome)
	req.Args = []string{"test", "--cold-cache", testDir}
	return nil
}
```
