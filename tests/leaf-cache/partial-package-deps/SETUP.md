# Scenario

**Feature**: partial leaf-cache across 3 leaves / 4 local packages

```
# packages
shared/{a,b,c} -> leaf-ab-1, leaf-ab-2
alone/d        -> leaf-d only

# product
leaf key includes local import closure sources
  -> edit alone/d invalidates leaf-d only
  -> edit shared/a invalidates leaf-ab-* only
```

## Preconditions

- Runtime parent builds binary + isolates `DOCTEST_LEAF_CACHE`.
- Fixture: 4 packages + 3 leaves under a temp module.

## Steps

1. Children call `preparePartialPackageDepsFixture`.
2. Run1 uses `-count=1` (seed PutPass, 0 Cached).
3. Mutate one package; Run2 without `-count` (warm peers).

## Context

- Proves Cached is per-leaf DAG, not whole-tree wipe.

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
	req.Timeout = 240 * time.Second
	// tests/leaf-cache -> tests -> module root
	modRoot := filepath.Clean(filepath.Join(d.DOCTEST_ROOT, "..", ".."))
	req.Bin = testbin.Ensure(t, modRoot)
	req.Env = isolateRuntimeEnv(t)
	return nil
}
```
