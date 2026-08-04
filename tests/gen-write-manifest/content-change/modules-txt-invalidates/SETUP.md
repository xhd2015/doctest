# Scenario

**Feature**: modules.txt-only change invalidates gomod-src cache and regenerates

```
first WriteGoMod (vendor with one module)
  -> seed tidy-done
change only vendor/modules.txt (+ package dir for new module)
second WriteGoMod
  -> gen go.mod includes new module
  -> doctest.tidy-done removed
```

## Steps

1. Prepare parent go.mod + vendor with one module.
2. First WriteGoMod; seed tidy-done.
3. Set ChangeSourceModulesTxt and create second package path for the new module.

```go
import (
	"os"
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Mode = "write-gomod-second"
	req.ModPath = "example.com/app"
	req.HasMod = true
	// Parent content-change may set ChangeSourceGoMod; this leaf is modules.txt-only.
	req.ChangeSourceGoMod = ""
	prepareFreshGen(t, req, "module example.com/app\n\ngo 1.21\n")

	vendor := filepath.Join(req.ModRoot, "vendor")
	depDir := filepath.Join(vendor, "example.com", "dep")
	if err := os.MkdirAll(depDir, 0755); err != nil {
		t.Fatalf("mkdir vendor dep: %v", err)
	}
	if err := os.WriteFile(filepath.Join(depDir, "dep.go"), []byte("package dep\n"), 0644); err != nil {
		t.Fatalf("write dep.go: %v", err)
	}
	firstModules := "# example.com/dep v1.0.0\n## explicit\nexample.com/dep\n"
	if err := os.WriteFile(filepath.Join(vendor, "modules.txt"), []byte(firstModules), 0644); err != nil {
		t.Fatalf("write modules.txt: %v", err)
	}

	firstWriteGoMod(t, req)
	seedTidyDone(t, req.GenDir)
	req.SnapGoModContentBefore = readFileOrEmpty(filepath.Join(req.GenDir, "go.mod"))

	otherDir := filepath.Join(vendor, "example.com", "other")
	if err := os.MkdirAll(otherDir, 0755); err != nil {
		t.Fatalf("mkdir other: %v", err)
	}
	if err := os.WriteFile(filepath.Join(otherDir, "other.go"), []byte("package other\n"), 0644); err != nil {
		t.Fatalf("write other.go: %v", err)
	}
	req.ChangeSourceModulesTxt = firstModules + "# example.com/other v2.0.0\n## explicit\nexample.com/other\n"
	return nil
}
```
