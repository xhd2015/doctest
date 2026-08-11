# Scenario

**Feature**: in-module multi-leaf tree imports parent `internal/greet`

```
# fixture: Hello + DefaultName, leaves leaf-a + leaf-b
createParentInternalMultiLeafModule
  -> ModuleRoot/tests with parent-internal imports
  -> GenDir under temp for layout inspect
```

## Preconditions

- Shared fixture helpers live on root (`createParentInternalMultiLeafModule`).
- Both leaves import `example.com/app/internal/greet` via root subject `Run`.

## Steps

1. Materialize parent module + two subject leaves.
2. Allocate inspectable `GenDir`.
3. Descendant leaves set cover flags or assert packaging.

```go
import (
	"os"
	"path/filepath"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	_ = d
	mod, tests := createParentInternalMultiLeafModule(t)
	req.ModuleRoot = mod
	req.TestDir = tests
	req.GenDir = filepath.Join(t.TempDir(), "gen")
	if err := os.MkdirAll(req.GenDir, 0755); err != nil {
		t.Fatalf("mkdir GenDir: %v", err)
	}
	return nil
}
```
