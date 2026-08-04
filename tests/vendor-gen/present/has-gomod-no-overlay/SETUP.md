# Scenario

**Feature**: modules that already have vendor go.mod get no overlay placeholder

```
# vendor has only example.com/dep with existing go.mod (no nogo without go.mod)
vendor/example.com/dep/go.mod present
  -> WriteGoMod
  -> replace example.com/dep => <modRoot>/vendor/example.com/dep
  -> NO genDir/vendor-gomod-overlay/example.com/dep/
  -> NO vendor-gomod-overlay.json (no placeholders needed)
```

## Steps

1. Start from present dirs, then rewrite vendor to a single module that already
   has go.mod (clears the default nogo entry so overlay should stay empty).
2. Run WriteGoMod; assert no overlay artifacts for that module and no overlay JSON.

## Context

- Complementary to `placeholder-gomod` (missing go.mod → overlay + JSON).
- Uses only-dep fixture so absence of overlay JSON is unambiguous.

```go
import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	if req.ModRoot == "" || req.VendorRoot == "" {
		t.Fatal("present Setup must set ModRoot and VendorRoot")
	}
	// Rewrite vendor: only sample dep with existing go.mod.
	req.NoGoModPath = ""
	req.NoGoModVersion = ""
	modulesTxt := strings.Join([]string{
		"# " + sampleDepPath + " " + sampleDepVersion,
		"## explicit; go 1.18",
		sampleDepPath,
		"",
	}, "\n")
	writeFile(t, filepath.Join(req.VendorRoot, "modules.txt"), modulesTxt)

	// Remove nogo tree from default present seed so no placeholder is needed.
	nogoDir := filepath.Join(req.VendorRoot, filepath.FromSlash(noGoModPath))
	if err := os.RemoveAll(nogoDir); err != nil {
		t.Fatalf("remove nogo fixture: %v", err)
	}

	depMod := filepath.Join(vendorModuleDir(req, req.SampleModPath), "go.mod")
	if !fileExists(depMod) {
		t.Fatalf("precondition: dep go.mod must exist at %s", depMod)
	}
	return nil
}
```
