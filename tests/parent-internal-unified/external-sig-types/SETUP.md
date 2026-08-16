# Scenario

**Feature**: Kind B expose compiles when product internal APIs use external package types

```
createParentInternalExternalSigModule
  model.Project + internal/rules.FixIgnore(model.Project, …)
  tests/leaf-a imports internal/rules
  -> runner.RunTest(tests, GenDir)
  -> Kind B expose for internal/rules must compile
  -> subject leaf PASS
```

Crime scene: generated expose body prints `model.Project` / `model.FixResult` in
the facade signature but only imports `srcpkg "…/internal/rules"` →
`undefined: model` → suite build failed.

## Preconditions

- Isolated temp module under `t.TempDir()` (parallel-safe).
- No cover flags required (failure is independent of coverage).

## Steps

1. Materialize product module with `model` + `internal/rules` + subject tree.
2. Allocate inspectable `GenDir`.
3. Run subject tree via root `Run` (no cover).

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
	mod, tests := createParentInternalExternalSigModule(t)
	req.ModuleRoot = mod
	req.TestDir = tests
	req.GenDir = filepath.Join(t.TempDir(), "gen")
	if err := os.MkdirAll(req.GenDir, 0755); err != nil {
		t.Fatalf("mkdir GenDir: %v", err)
	}
	req.WithCover = false
	req.CoverPath = ""
	req.CoverPkg = ""
	req.CoverMode = ""
	return nil
}
```
