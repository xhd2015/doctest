# Scenario

**Feature**: the doctest module root is two levels above this test tree (`DOCTEST_ROOT/../..`)

```
# maps doc-style leaves to Go test packages
doctest build/test -> mapping-gen -> leaf <-> Go package

# gen-dir modes
auto gen-dir -> one package per leaf | explicit gen-dir -> user-specified layout
```

## Preconditions
- The doctest module root is two levels above this test tree (`DOCTEST_ROOT/../..`).
- The mapping-gen feature mirrors source directory structure under the cache root.
- Each leaf sets up a temporary project and verifies generated output.

## Steps
1. Build the doctest binary from the module root.
2. Execute the binary given by `req.Bin`.
3. Capture stdout, stderr, exit code, and the raw execution error.

## Context
- Each test creates a temp Go project with a doctest tree.
- The generated output is inspected for per-leaf directory structure.
- Shared go.mod is at project root level in the generated tree.

```go
import (
"github.com/xhd2015/doctest/session"
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
	"github.com/xhd2015/doctest/libdoc/testbin"
	"github.com/xhd2015/doctest/libdoc/testtree"
)

var bt = "`" + "`" + "`"
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Timeout = 120 * time.Second

	req.Bin = testbin.Ensure(t, filepath.Join(d.DOCTEST_ROOT, "..", ".."))
	return nil
}
func doctestGoBlock(code string) string {
	return "## Test\n\n" + bt + "go\n" + code + bt + "\n"
}
func doctestBody(runCode string) string {
	return "import \"testing\"\n\ntype Request struct{ Name string }\ntype Response struct{ Message string }\n\n" + runCode
}
func rootSetupContent() string {
	return doctestGoBlock("import \"testing\"\nfunc Setup(t *testing.T, req *Request) error { _ = req; return nil }")
}
func leafSetupContent() string {
	return doctestGoBlock("import \"testing\"\nfunc Setup(t *testing.T, req *Request) error { _ = req; return nil }")
}
func leafAssertContent() string {
	return doctestGoBlock("import \"testing\"\nfunc Assert(t *testing.T, req *Request, resp *Response, err error) {}")
}
func createDoctestRoot(dir string, runCode string) error {
	if err := os.WriteFile(filepath.Join(dir, "DOCTEST.md"), []byte(testtree.MinimalDOCTEST(doctestBody(runCode))), 0644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, "SETUP.md"), []byte(rootSetupContent()), 0644); err != nil {
		return err
	}
	return nil
}
func createDoctestLeaf(dir string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, "SETUP.md"), []byte(leafSetupContent()), 0644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, "ASSERT.md"), []byte(leafAssertContent()), 0644); err != nil {
		return err
	}
	return nil
}
func createTempProject(t *testing.T, dirName string) string {
	t.Helper()
	tmp := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmp, "go.mod"), []byte("module testproj\ngo 1.21\n"), 0644); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}
	testDir := filepath.Join(tmp, dirName)
	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatalf("mkdir test dir: %v", err)
	}
	runCode := "func Run(t *testing.T, req *Request) (*Response, error) { return &Response{}, nil }"
	if err := createDoctestRoot(testDir, runCode); err != nil {
		t.Fatalf("create doctest root: %v", err)
	}
	return testDir
}
func assertFileExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected file %s to exist: %v", path, err)
	}
}
func assertFileNotExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); err == nil {
		t.Fatalf("expected file %s to not exist", path)
	}
}
func parseGenDir(stderr string) string {
	idx := strings.Index(stderr, "→ ")
	if idx < 0 {
		return ""
	}
	rest := stderr[idx+len("→ "):]
	end := strings.IndexFunc(rest, func(r rune) bool { return r == '\n' || r == ' ' })
	if end < 0 {
		return strings.TrimSpace(rest)
	}
	return strings.TrimSpace(rest[:end])
}
func findGoModDir(startDir string) string {
	dir := startDir
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}
```
