---
label: heavy
---

## Expected

- Both `doctest test` runs exit 0.
- Cached `assert.go` and `go.mod` have identical MD5 digests after second run.
- Modification times are unchanged (write-once semantics).

## Exit Code

- Exit code 0.

```go
import (
	"path/filepath"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("doctest subprocess failed: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected second-run idempotent test to pass, got exit %d\nstdout:\n%s\nstderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	cacheDir := expectedAssertCacheDir(t)
	afterAssertMtime, afterAssertDigest := snapshotFileState(t, filepath.Join(cacheDir, "assert.go"))
	afterGoModMtime, afterGoModDigest := snapshotFileState(t, filepath.Join(cacheDir, "go.mod"))
	if afterAssertDigest != beforeAssertDigest {
		t.Fatalf("assert.go content changed on second run")
	}
	if afterGoModDigest != beforeGoModDigest {
		t.Fatalf("go.mod content changed on second run")
	}
	if afterAssertMtime != beforeAssertMtime {
		t.Fatalf("assert.go mtime changed on second run: before=%d after=%d", beforeAssertMtime, afterAssertMtime)
	}
	if afterGoModMtime != beforeGoModMtime {
		t.Fatalf("go.mod mtime changed on second run: before=%d after=%d", beforeGoModMtime, afterGoModMtime)
	}
}
```