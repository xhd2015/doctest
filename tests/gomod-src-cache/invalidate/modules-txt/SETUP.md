# Scenario

**Feature**: modules.txt-only change invalidates gomod-src and regenerates vendor inject

```
first write (vendor with one module) -> seed tidy-done
change only vendor/modules.txt (+ package dir for new module)
second write
  -> gen go.mod includes new module
  -> tidy-done removed
```

## Steps

1. Parent go.mod + vendor with example.com/dep; first write; seed tidy-done.
2. Add example.com/other package dir; set ChangeSourceModulesTxt.

```go
import (
	"os"
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Mode = "write-second"
	req.ModPath = defaultModPath
	req.HasMod = true
	prepareFreshGen(t, req, "module example.com/app\n\ngo 1.21\n")

	vendor := filepath.Join(req.ModRoot, "vendor")
	depDir := filepath.Join(vendor, "example.com", "dep")
	if err := os.MkdirAll(depDir, 0755); err != nil {
		t.Fatalf("mkdir vendor dep: %v", err)
	}
	if err := os.WriteFile(filepath.Join(depDir, "dep.go"), []byte("package dep\n"), 0644); err != nil {
		t.Fatalf("write dep.go: %v", err)
	}
	// dep has go.mod so no placeholder required for first module.
	if err := os.WriteFile(filepath.Join(depDir, "go.mod"), []byte("module example.com/dep\n\ngo 1.18\n"), 0644); err != nil {
		t.Fatalf("write dep go.mod: %v", err)
	}
	firstModules := "# example.com/dep v1.0.0\n## explicit\nexample.com/dep\n"
	if err := os.WriteFile(filepath.Join(vendor, "modules.txt"), []byte(firstModules), 0644); err != nil {
		t.Fatalf("write modules.txt: %v", err)
	}

	firstWrite(t, req)
	seedTidyDone(t, req.GenDir)
	req.SnapGoModContentBefore = readFileOrEmpty(filepath.Join(req.GenDir, "go.mod"))

	otherDir := filepath.Join(vendor, "example.com", "other")
	if err := os.MkdirAll(otherDir, 0755); err != nil {
		t.Fatalf("mkdir other: %v", err)
	}
	if err := os.WriteFile(filepath.Join(otherDir, "other.go"), []byte("package other\n"), 0644); err != nil {
		t.Fatalf("write other.go: %v", err)
	}
	if err := os.WriteFile(filepath.Join(otherDir, "go.mod"), []byte("module example.com/other\n\ngo 1.18\n"), 0644); err != nil {
		t.Fatalf("write other go.mod: %v", err)
	}
	req.ChangeSourceModulesTxt = firstModules + "# example.com/other v2.0.0\n## explicit\nexample.com/other\n"
	return nil
}
```
