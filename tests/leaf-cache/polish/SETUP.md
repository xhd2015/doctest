# Scenario

**Feature**: P3 polish — selective invalidation, tree isolation, help docs

```
# selective: edit one leaf, sibling stays Cached
# local dep: edit imported package, leaf re-runs
# isolation: warm treeA does not warm treeB
# docs: test --help lists leaf-cache flags
```

## Preconditions

- Nested CLI via testbin (same as runtime/).
- Isolated DOCTEST_LEAF_CACHE + fresh GOCACHE per run.
- Leaves labeled `heavy` where nested compile is required.

## Steps

1. Ensure binary; isolate env.
2. Child selects fixture and Args multi-run sequence.

## Context

- Extends P2 warm-skip with multi-leaf and multi-tree scenarios.

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
	req.Timeout = 240 * time.Second
	modRoot := filepath.Clean(filepath.Join(d.DOCTEST_ROOT, "..", ".."))
	req.Bin = testbin.Ensure(t, modRoot)
	req.Env = isolateRuntimeEnv(t)
	return nil
}
```
