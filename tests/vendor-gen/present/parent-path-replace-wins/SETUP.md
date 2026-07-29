# Scenario

**Feature**: parent filesystem path replace for module M wins over vendor inject

```
parent go.mod:
  replace example.com/dep => ./work_around/example.com/dep
+ vendor/modules.txt also lists example.com/dep (and nogo)
  -> WriteGoMod
  -> gen go.mod: exactly one replace for example.com/dep
     RHS = abs work_around path (NOT …/vendor/example.com/dep)
  -> example.com/nogo still gets vendor require+replace
```

## Steps

1. Inherit present vendor fixture (`example.com/dep` + `example.com/nogo`).
2. Create a local work_around tree for the same module path as SampleModPath.
3. Rewrite parent go.mod with a filesystem path replace for SampleModPath → work_around.
4. Leave vendor/modules.txt + vendor tree for that module in place (conflict condition).

## Context

- Distinct from `coexists-parent`: there the parent path replace is for a module
  **not** listed in modules.txt; here the **same** module appears in both parent
  replace and vendor modules.txt.
- Classic TDD: **RED** today — `vendorBridgeForModRoot` always injects
  `replace M => …/vendor/M` even when `readExtraReplaces` already copied the
  parent path replace for M (dual replace → `go mod tidy` conflict).

```go
import (
	"os"
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	if req.SampleModPath == "" || req.NoGoModPath == "" {
		t.Fatal("present Setup must set SampleModPath and NoGoModPath")
	}
	// Parent path replace for the SAME module as a vendor entry.
	// Mirrors real projects (e.g. replace github.com/gansidui/go-utils => ./tools/work_around/…).
	req.ParentLocalRel = "./work_around/" + req.SampleModPath
	workAroundAbs := filepath.Join(req.ModRoot, "work_around", filepath.FromSlash(req.SampleModPath))
	req.ParentLocalAbs = workAroundAbs
	if err := os.MkdirAll(workAroundAbs, 0755); err != nil {
		t.Fatalf("mkdir work_around: %v", err)
	}
	writeFile(t, filepath.Join(workAroundAbs, "go.mod"),
		"module "+req.SampleModPath+"\n\ngo 1.19\n")
	writeFile(t, filepath.Join(workAroundAbs, "dep.go"),
		"package dep\n\n// PARENT_PATH_REPLACE_MARKER\nconst FromParent = true\n")

	body := "module " + req.ModPath + "\n\ngo 1.19\n\n" +
		"replace " + req.SampleModPath + " => " + req.ParentLocalRel + "\n"
	seedParentGoMod(t, req, body)

	// Vendor fixture must still list SampleModPath (overlap with parent replace).
	if !fileExists(filepath.Join(req.VendorRoot, "modules.txt")) {
		t.Fatal("vendor/modules.txt must remain after parent go.mod rewrite")
	}
	if !fileExists(vendorModuleDir(req, req.SampleModPath)) {
		t.Fatal("vendor tree for SampleModPath must remain (conflict condition)")
	}
	if !fileExists(vendorModuleDir(req, req.NoGoModPath)) {
		t.Fatal("vendor tree for NoGoModPath must remain (non-overridden module)")
	}
	return nil
}
```
