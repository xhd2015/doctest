---
label: heavy
---

## Expected
- Command succeeds (exit code 0).
- Generated output has per-leaf directory structure with `_test.go` in each leaf dir.
- `go build ./...` runs successfully from the go.mod level in the generated tree.
- Test files contain compile-only stubs (compiles but doesn't execute).
- No `doctest.hash` file exists.

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
		t.Fatalf("build failed: %v\nstderr:\n%s", err, resp.Stderr)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected exit 0, got %d\nstderr:\n%s", resp.ExitCode, resp.Stderr)
	}

	if req.GenDir == "" {
		t.Fatal("req.GenDir is empty; grouping Setup must set request-local gen dir")
	}

	// req-local path, Parallel-safe
	case1 := filepath.Join(req.GenDir, "tests", "feature", "case1")
	case2 := filepath.Join(req.GenDir, "tests", "feature", "case2")

	// Unified mode: non-test leaf packages (leaf.go), not classic *_test.go.
	assertFileExists(t, filepath.Join(case1, "leaf.go"))
	assertFileExists(t, filepath.Join(case2, "leaf.go"))

	// Verify compile-only stubs: generated files should contain compileOnly usage
	data, err := os.ReadFile(filepath.Join(case1, "leaf.go"))
	if err != nil {
		t.Fatalf("read leaf.go: %v", err)
	}
	if !strings.Contains(string(data), "compileOnly") {
		t.Fatalf("expected compileOnly in generated build test source:\n%s", string(data))
	}

	// Shared go.mod at project root
	assertFileExists(t, filepath.Join(req.GenDir, "go.mod"))

	// No hash file
	filepath.Walk(req.GenDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(info.Name(), ".hash") {
			t.Fatalf("found unexpected hash file: %s", path)
		}
		return nil
	})
}
```
