# Scenario

**Feature**: missing vendor go.mod gets xgo-style overlay placeholder only

```
vendor/example.com/nogo/   # package only, no go.mod
  -> WriteGoMod
  -> project vendor still has no go.mod (read-only)
  -> genDir/vendor-gomod-overlay/example.com/nogo/go.mod
       module example.com/nogo
       go 1.17   # from modules.txt ## go metadata when present
  -> genDir/vendor-gomod-overlay.json
       Replace[abs(project vendor/.../nogo/go.mod)] = abs(placeholder)
  -> replace example.com/nogo => <modRoot>/vendor/example.com/nogo
  -> NO package files under vendor-gomod-overlay/<mod>/ (no hardlink/copy)
```

## Steps

1. Inherit present fixture (nogo has package source, no go.mod).
2. Confirm pre-Run that nogo lacks go.mod.
3. After Run, assert overlay placeholder, overlay JSON mapping, project
   immutability, and replace → project vendor.

## Context

- Modules that already have go.mod (`example.com/dep`) are covered by sibling
  `has-gomod-no-overlay`; this leaf only asserts the missing-go.mod case.
- xgo-style: `module <path>` + `go <ver>`; go version from modules.txt `## … go X.Y`
  or parent go directive when metadata omits it.

```go
import (
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	if req.NoGoModPath == "" {
		t.Fatal("NoGoModPath required")
	}
	nogoMod := filepath.Join(vendorModuleDir(req, req.NoGoModPath), "go.mod")
	if fileExists(nogoMod) {
		t.Fatalf("precondition: %s must not exist before WriteGoMod", nogoMod)
	}
	// dep should already have go.mod from seed
	depMod := filepath.Join(vendorModuleDir(req, req.SampleModPath), "go.mod")
	if !fileExists(depMod) {
		t.Fatalf("precondition: dep go.mod should exist at %s", depMod)
	}
	return nil
}
```
