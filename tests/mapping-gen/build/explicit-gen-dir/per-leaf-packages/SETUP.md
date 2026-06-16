# Scenario

**Feature**: a --gen-dir is specified for build mode

```
# maps doc-style leaves to Go test packages
doctest build/test -> mapping-gen -> leaf <-> Go package

# gen-dir modes
auto gen-dir -> one package per leaf | explicit gen-dir -> user-specified layout
```

## Preconditions
- A --gen-dir is specified for build mode.
- A temp project with 2 leaves under a grouping directory exists.

## Steps
1. Create a project with `go.mod`, doctest root at `tests/`.
2. Create a grouping directory `tests/feature/` with a SETUP.md.
3. Create two leaves: `tests/feature/case1/` and `tests/feature/case2/`.
4. Run `doctest build <test-dir> --gen-dir <genDir> -v`.

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

	featureDir := filepath.Join(testDir, "feature")
	if err := os.MkdirAll(featureDir, 0755); err != nil {
		t.Fatalf("mkdir feature dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(featureDir, "SETUP.md"), []byte(doctestGoBlock("import \"testing\"\nfunc Setup(t *testing.T, req *Request) error { _ = req; return nil }")), 0644); err != nil {
		t.Fatalf("write feature setup: %v", err)
	}

	if err := createDoctestLeaf(filepath.Join(featureDir, "case1")); err != nil {
		t.Fatalf("create case1: %v", err)
	}
	if err := createDoctestLeaf(filepath.Join(featureDir, "case2")); err != nil {
		t.Fatalf("create case2: %v", err)
	}

	req.Args = append(req.Args, "-v", testDir)
	return nil
}
```
