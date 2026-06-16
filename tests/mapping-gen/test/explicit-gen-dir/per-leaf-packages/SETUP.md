# Scenario

**Feature**: a temp project with 2 leaves under a nested grouping directory exists

```
# maps doc-style leaves to Go test packages
doctest build/test -> mapping-gen -> leaf <-> Go package

# gen-dir modes
auto gen-dir -> one package per leaf | explicit gen-dir -> user-specified layout
```

## Preconditions
- A temp project with 2 leaves under a nested grouping directory exists.

## Steps
1. Create a project with `go.mod`, doctest root at `tests/`.
2. Create a grouping directory `tests/category/` with a SETUP.md (no ASSERT.md).
3. Create two leaves: `tests/category/leaf-a/` and `tests/category/leaf-b/`.
4. Run `doctest test <test-dir> --gen-dir <genDir> -v`.

```go
import (
	"os"
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	proj := t.TempDir()
	if err := os.WriteFile(filepath.Join(proj, "go.mod"), []byte("module testproj\ngo 1.21\n"), 0644); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}

	testDir := filepath.Join(proj, "tests")
	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatalf("mkdir test dir: %v", err)
	}
	runCode := "func Run(t *testing.T, req *Request) (*Response, error) { return &Response{}, nil }"
	if err := createDoctestRoot(testDir, runCode); err != nil {
		t.Fatalf("create doctest root: %v", err)
	}

	groupDir := filepath.Join(testDir, "category")
	if err := os.MkdirAll(groupDir, 0755); err != nil {
		t.Fatalf("mkdir group dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(groupDir, "SETUP.md"), []byte(doctestGoBlock("import \"testing\"\nfunc Setup(t *testing.T, req *Request) error { _ = req; return nil }")), 0644); err != nil {
		t.Fatalf("write group setup: %v", err)
	}

	if err := createDoctestLeaf(filepath.Join(groupDir, "leaf_a")); err != nil {
		t.Fatalf("create leaf a: %v", err)
	}
	if err := createDoctestLeaf(filepath.Join(groupDir, "leaf_b")); err != nil {
		t.Fatalf("create leaf b: %v", err)
	}

	req.Args = append(req.Args, "-v", testDir)
	return nil
}
```
