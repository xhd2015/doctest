# Scenario

**Bug**: parent module→module replace for module M duals with vendor inject

```
parent go.mod:
  replace example.com/dep => example.com/dep v1.3.0
  # (same-path version pin; also covers fork form: => example.com/other v1.0.0)
+ vendor/modules.txt also lists example.com/dep (and nogo)
  -> WriteGoMod
  -> gen go.mod: exactly one replace for example.com/dep
     RHS = parent module→module target (NOT …/vendor/example.com/dep)
  -> example.com/nogo still gets vendor require+replace
```

## Steps

1. Inherit present vendor fixture (`example.com/dep` + `example.com/nogo`).
2. Rewrite parent go.mod with a **module→module** replace for SampleModPath
   (version pin, not a filesystem path).
3. Leave vendor/modules.txt + vendor tree for that module in place (conflict
   condition — same as real monorepos with `replace github.com/gogo/protobuf =>
   github.com/gogo/protobuf v1.3.0` plus vendored protobuf).

## Context

- Distinct from `parent-path-replace-wins`: there the parent replace RHS is a
  **filesystem path** (`./work_around/…`); here RHS is another **module path +
  version** (or same path with pinned version).
- Distinct from `coexists-parent`: there the parent path replace is for a module
  **not** listed in modules.txt; here the **same** module appears in both parent
  replace and vendor modules.txt.
- Classic TDD: **RED** today — `readExtraReplaces` copies module→module replaces
  into gen go.mod but does **not** put them in the skip set used by
  `vendorBridgeForModRoot` (only filesystem path replaces are skipped), so M
  gets both parent replace and `replace M => …/vendor/M` (dual replace →
  `go mod tidy` conflict).

```go
import (
	"path/filepath"
	"testing"
)

// parentModuleReplaceVersion is a distinctive pin so Assert can see the parent
// module→module RHS preserved (modules.txt still lists SampleModVersion v1.2.3).
const parentModuleReplaceVersion = "v1.3.0"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	if req.SampleModPath == "" || req.NoGoModPath == "" {
		t.Fatal("present Setup must set SampleModPath and NoGoModPath")
	}
	// Module→module replace for the SAME module as a vendor entry.
	// Mirrors pricing-center / gogo: replace M => M vX.Y.Z (version pin).
	// Store expected RHS tokens for Assert (reuse ParentLocalRel/Abs slots as
	// "module path" and "version" labels — no filesystem work_around tree).
	req.ParentLocalRel = req.SampleModPath // RHS module path (same path pin)
	req.ParentLocalAbs = parentModuleReplaceVersion

	body := "module " + req.ModPath + "\n\ngo 1.19\n\n" +
		"replace " + req.SampleModPath + " => " + req.SampleModPath + " " + parentModuleReplaceVersion + "\n"
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
