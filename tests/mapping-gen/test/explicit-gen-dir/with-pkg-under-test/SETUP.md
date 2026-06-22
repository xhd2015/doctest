# Scenario

**Feature**: a temp project with a package under test declaration in the root SETUP.md

```
# maps doc-style leaves to Go test packages
doctest build/test -> mapping-gen -> leaf <-> Go package

# gen-dir modes
auto gen-dir -> one package per leaf | explicit gen-dir -> user-specified layout
```

## Preconditions
- A temp project with a package under test declaration in the root SETUP.md.

## Steps
1. Create a project with `go.mod` and a Go source file at `src/calc.go`.
2. The doctest root SETUP.md declares `- Package under test: calc`.
3. Create 2 leaves.
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

	srcDir := filepath.Join(proj, "src")
	if err := os.MkdirAll(srcDir, 0755); err != nil {
		t.Fatalf("mkdir src dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(srcDir, "calc.go"), []byte("package calc\n\nfunc Add(a, b int) int { return a + b }\n"), 0644); err != nil {
		t.Fatalf("write calc.go: %v", err)
	}

	testDir := filepath.Join(proj, "tests")
	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatalf("mkdir test dir: %v", err)
	}
	runCode := "func Run(t *testing.T, req *Request) (*Response, error) { return &Response{}, nil }"
	if err := createDoctestRoot(testDir, runCode); err != nil {
		t.Fatalf("create doctest root: %v", err)
	}
	rootSetup := "- Package under test: calc\n\n" + rootSetupContent()
	if err := os.WriteFile(filepath.Join(testDir, "SETUP.md"), []byte(rootSetup), 0644); err != nil {
		t.Fatalf("write root SETUP.md: %v", err)
	}

	for _, name := range []string{"leaf_a", "leaf_b"} {
		if err := createDoctestLeaf(filepath.Join(testDir, name)); err != nil {
			t.Fatalf("create leaf %s: %v", name, err)
		}
	}

	req.Args = append(req.Args, "-v", testDir)
	return nil
}
```
