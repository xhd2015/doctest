# Scenario

**Feature**: doctest test wires leaf-cache into skip + summary Cached

```
# isolated env
DOCTEST_LEAF_CACHE + DOCTEST_CACHE_HOME; fresh GOCACHE per run
  -> doctest test <fixture>
  -> pass -> PutPass(key)
  -> warm skip -> Cached in summary
```

## Preconditions

- Nested CLI integration — single-tree leaf-cache product path (**GREEN**).
- Requires a freshly built doctest binary (`testbin.Ensure`).
- Fixture trees are temp mini doctest projects (not this tree's key/store leaves).
- Fresh `GOCACHE` per invocation so `N Cached` means programmatic leaf-cache skips only.

## Steps

1. Resolve module root from `d.DOCTEST_ROOT` and ensure binary.
2. Set `Op=runtime_multi`, timeout, and isolateRuntimeEnv.
3. Child leaves choose pass/fail fixture and second-run flags.

## Context

- Significance: single-tree product surface; workspace `./...` is under
  `workspace/**`, multi-arg under `cli-plan/**` — same Cached/disable policy.
- Labels: leaves use `heavy` (nested compile + two runs).

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
	req.Op = "runtime_multi"
	req.Timeout = 180 * time.Second
	// tests/leaf-cache -> tests -> module root
	modRoot := filepath.Clean(filepath.Join(d.DOCTEST_ROOT, "..", ".."))
	req.Bin = testbin.Ensure(t, modRoot)
	req.Env = isolateRuntimeEnv(t)
	return nil
}
```
