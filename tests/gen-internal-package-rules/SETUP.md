# Scenario

**Feature**: shared product doctest binary for gen-internal RED contracts

```
testbin.Ensure(moduleRoot) -> req.Bin
leaf Setup materializes fixture -> doctest test <subject>
```

## Preconditions

- Module root is two levels above this tree (`DOCTEST_ROOT/../..`).
- `go` on PATH.

## Steps

1. Build/share doctest binary via `testbin.Ensure`.
2. Leaves copy fixtures and set `Args` / `WorkDir` / `Kind`.

```go
import (
	"path/filepath"
	"testing"
	"time"

	"github.com/xhd2015/doctest/libdoc/testbin"
	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	req.Timeout = 180 * time.Second
	req.Bin = testbin.Ensure(t, filepath.Join(d.DOCTEST_ROOT, "..", ".."))
	req.Args = nil
	req.Env = nil
	req.WorkDir = ""
	req.Kind = ""
	return nil
}
```
