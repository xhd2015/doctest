# Scenario

**Feature**: with vendor/ + modules.txt, WriteGoMod injects vendor bridge

```
# project has go.mod + vendor/modules.txt + vendored module trees
modRoot/vendor/modules.txt
  -> WriteGoMod
  -> gen go.mod: require + replace for each modules.txt module
  -> placeholder go.mod under vendored modules missing go.mod
```

## Preconditions

- `modRoot` has both `go.mod` and `vendor/` with a tiny `modules.txt`.
- Default fixture modules: `example.com/dep` (has go.mod + marker source) and
  `example.com/nogo` (package only, no go.mod).
- Classic TDD: leaves under this branch are **RED** until vendor-aware WriteGoMod.

## Steps

1. Create temp modRoot + genDir.
2. Seed parent go.mod (default go 1.19 unless leaf overrides).
3. Seed tiny vendor tree via `seedTinyVendor`.

## Context

- Leaves assert different facets of the same fixture; they do not re-seed
  conflicting vendor layouts unless they override after this Setup.

```go
import (
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	prepareFreshDirs(t, req)
	seedParentGoMod(t, req, "module "+req.ModPath+"\n\ngo 1.19\n")
	seedTinyVendor(t, req)
	if !fileExists(filepath.Join(req.VendorRoot, "modules.txt")) {
		t.Fatal("present Setup must create vendor/modules.txt")
	}
	return nil
}
```
