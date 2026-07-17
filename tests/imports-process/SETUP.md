# Scenario

**Feature**: the doctest module root is two levels above this test tree (`DOCTEST_ROOT/../..`)

```
# Go import processing during code generation
doctest build -> parse imports -> remove unused -> report syntax errors
```

## Preconditions
- The doctest module root is two levels above this test tree (`DOCTEST_ROOT/../..`).
- Each leaf creates a temporary doctest project and runs `doctest test` on it.
- The doctest binary is built fresh for each test.

## Steps
1. Build the doctest binary from the module root.
2. Create a temp Go project with a doctest tree containing specific imports.
3. Run `doctest test <test-dir>` and capture output.

## Context
- These tests verify the `imports.Process` behavior in `WriteGeneratedCase`.
- The generated test code must compile and pass for the unused-import case.
- The syntax error case must produce a clear error message.

```go
import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
	"github.com/xhd2015/doctest/libdoc/testbin"
	"github.com/xhd2015/doctest/libdoc/testtree"
)

var bt = "\x60\x60\x60"
func Setup(t *testing.T, req *Request) error {
	req.Timeout = 120 * time.Second

	req.Bin = testbin.Ensure(t, filepath.Join(DOCTEST_ROOT, "..", ".."))
	return nil
}
func doctestGoBlock(code string) string {
	return "\n## Test\n\n" + bt + "go\n" + code + "\n" + bt + "\n"
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
func createDoctestLeaf(dir string, setupContent string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, "SETUP.md"), []byte(setupContent), 0644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, "ASSERT.md"), []byte(leafAssertContent()), 0644); err != nil {
		return err
	}
	return nil
}
```
