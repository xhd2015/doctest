---
label: heavy
---

## Expected
- Command succeeds (exit code 0).
- Generated output has per-leaf directory structure: each leaf gets its own dir with `_test.go`.
- A shared `go.mod` exists at the project root level in the generated tree.
- No `doctest.hash` file exists anywhere in the generated tree.
- Two test cases are discovered and both pass.

## Exit Code
- Exit code 0.

```go
import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("run failed: %v\nstderr:\n%s", err, resp.Stderr)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected exit 0, got %d\nstderr:\n%s", resp.ExitCode, resp.Stderr)
	}

	// Verify per-leaf directory structure
	leafA := filepath.Join(genDir, "tests", "category", "leaf_a")
	leafB := filepath.Join(genDir, "tests", "category", "leaf_b")

	// Unified mode: non-test leaf packages (leaf.go), not classic *_test.go.
	assertFileExists(t, filepath.Join(leafA, "leaf.go"))
	assertFileExists(t, filepath.Join(leafB, "leaf.go"))

	// Each leaf dir should NOT have its own go.mod — go.mod is shared at project root
	assertFileNotExists(t, filepath.Join(leafA, "go.mod"))
	assertFileNotExists(t, filepath.Join(leafB, "go.mod"))

	// Shared go.mod at project root level in generated tree
	assertFileExists(t, filepath.Join(genDir, "go.mod"))

	// No doctest.hash file anywhere in genDir
	filepath.Walk(genDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(info.Name(), ".hash") {
			t.Fatalf("found unexpected hash file: %s", path)
		}
		return nil
	})

	// Two test cases ran
	if !strings.Contains(resp.Stdout, "leaf_a") || !strings.Contains(resp.Stdout, "leaf_b") {
		t.Fatalf("expected leaf_a and leaf_b in test output:\n%s", resp.Stdout)
	}
}
```
