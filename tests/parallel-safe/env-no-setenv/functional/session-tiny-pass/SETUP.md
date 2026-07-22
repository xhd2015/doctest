# Scenario

**Feature**: `doctest test` on a tiny fixture exits 0 with nested session id inject

```
# session smoke
temp module + mytest/ one-leaf tree
  leaf Assert: d.DOCTEST_SESSION_ID != ""
doctest test mytest/
  -> exit 0
  -> child suite got SESSION_ID via go test cmd.Env (not parent process Setenv after P1)
```

## Preconditions

- Nested leaf Assert requires non-empty `d.DOCTEST_SESSION_ID`.
- Deeper Once/replace coverage remains in `tests/session-inject/`.

## Steps

1. Create temp module with session-checking leaf.
2. Run `doctest test <testDir>`.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	_, testDir := createTempModuleFixture(t, true)
	req.FixtureDir = testDir
	req.Args = []string{"test", testDir}
	return nil
}
```
