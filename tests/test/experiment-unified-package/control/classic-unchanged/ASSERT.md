## Expected

- Suite run succeeds (`RunErr` empty) with unified flag off.
- At least two leaf `*_test.go` files exist under gen paths for `a`/`b`.
- Unified-only dirs are **not** required: no requirement that `__registry` /
  `__allleaves` / suite-only packaging appear.
- Displayed `go test` package args are **not** a single `suite` package only
  (classic multi-leaf: at least two packages, or packages that are not solely suite).

```go
import (
	"path/filepath"
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected harness error: %v", err)
	}
	if resp.RunErr != "" {
		t.Fatalf("classic flag-off RunTest failed: %s\nstdout:\n%s\nstderr:\n%s",
			resp.RunErr, resp.Stdout, resp.Stderr)
	}
	if resp.GenDir == "" {
		t.Fatal("GenDir empty; cannot inspect classic layout")
	}
	if len(resp.LeafTestGoFiles) < 2 {
		t.Fatalf("classic expects ≥2 leaf *_test.go under a/b, got %v (all go=%v) gen=%s",
			basenames(resp.LeafTestGoFiles), basenames(resp.GoFiles), resp.GenDir)
	}
	// Must not look like unified suite-only packaging.
	pkgs := resp.GoTestPackageArgs
	if len(pkgs) == 1 && strings.Contains(pkgs[0], "suite") {
		t.Fatalf("classic flag-off must not use suite-only go test; pkgs=%v line=%q",
			pkgs, resp.GoTestDisplayLine)
	}
	// Multi-package classic shape is ideal; if display line is sparse, ≥2 leaf
	// *_test.go files already prove classic multi-leaf packaging.
	if len(pkgs) >= 2 {
		return
	}
}

func basenames(paths []string) []string {
	out := make([]string, len(paths))
	for i, p := range paths {
		out[i] = filepath.Base(p)
	}
	return out
}
```
