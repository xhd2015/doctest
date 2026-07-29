# Scenario

**Feature**: placeholder go.mod is created for vendored modules that lack one

```
vendor/example.com/nogo/   # package only, no go.mod
  -> WriteGoMod
  -> vendor/example.com/nogo/go.mod exists
       module example.com/nogo
       go 1.17   # from modules.txt ## go metadata when present
```

## Steps

1. Inherit present fixture (nogo has package source, no go.mod).
2. Confirm pre-Run that nogo lacks go.mod.
3. After Run, assert placeholder go.mod content.

## Context

- Modules that already have go.mod (`example.com/dep`) must not be required to
  lose their existing module file; this leaf only asserts the missing case.
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
