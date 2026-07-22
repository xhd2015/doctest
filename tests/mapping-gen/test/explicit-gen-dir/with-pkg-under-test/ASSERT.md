---
label: heavy
---

## Expected
- Command succeeds (exit code 0).
- Source files from the package under test are copied to each leaf directory.
- Copied files have renamed package declarations (e.g. `package calc_tc`).

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

	if req.GenDir == "" {
		t.Fatal("req.GenDir is empty; grouping Setup must set request-local gen dir")
	}

	// Source files should be copied to each leaf directory (req-local path, Parallel-safe)
	leafA := filepath.Join(req.GenDir, "tests", "leaf_a")
	leafB := filepath.Join(req.GenDir, "tests", "leaf_b")

	for _, leafDir := range []string{leafA, leafB} {
		calcPath := filepath.Join(leafDir, "calc.go")
		assertFileExists(t, calcPath)

		data, err := os.ReadFile(calcPath)
		if err != nil {
			t.Fatalf("read %s: %v", calcPath, err)
		}
		if !strings.Contains(string(data), "package calc_tc") {
			t.Fatalf("expected package calc_tc in %s, got:\n%s", calcPath, string(data))
		}
	}

	// Source files should NOT be in the grouping dir or root
	assertFileNotExists(t, filepath.Join(req.GenDir, "tests", "calc.go"))
	assertFileNotExists(t, filepath.Join(req.GenDir, "calc.go"))
}
```
