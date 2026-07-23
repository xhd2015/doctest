# Scenario

**Feature**: default in-process CLI; optional product binary for true e2e

```
# default L2: in-process CLI dispatch
cli.RunWithWriter(stdout, req.Args) -> stdout/stderr/exit

# optional L3: real product binary (UseCLI)
testbin.Ensure(moduleRoot) -> req.Bin
req.Bin <args> -> subprocess capture
```

## Preconditions

- Module root is the parent of this test tree (`DOCTEST_ROOT/..`).
- Default leaves do **not** require a product binary.
- Leaves that set `UseCLI` build/share a binary via `testbin.Ensure` and must
  carry `label: e2e`.

## Steps

1. Set a generous default timeout for nested selftests.
2. If `UseCLI` is already set (e2e leaf Setup), ensure `req.Bin`.
3. Otherwise leave Bin empty — `Run` uses in-process CLI.

## Context

- Agent e2e trees (implementer, etc.) set `UseCLI` and may need `fake-codex` on PATH.

```go
import (
	"path/filepath"
	"testing"
	"time"

	"github.com/xhd2015/doctest/session"
	"github.com/xhd2015/doctest/libdoc/testbin"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	// Nested doctest test + go test can exceed 30s on contended runners.
	if req.Timeout <= 0 {
		req.Timeout = 2 * time.Minute
	}
	if req.UseCLI {
		req.Bin = testbin.Ensure(t, filepath.Join(d.DOCTEST_ROOT, ".."))
	}
	return nil
}
```
