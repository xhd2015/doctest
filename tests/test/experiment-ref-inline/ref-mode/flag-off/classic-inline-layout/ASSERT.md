## Expected

- Suite run succeeds (`RunErr` empty).
- At least two generated `*.go` files define `func ExperimentP1RootMarker`
  (classic: one copy per leaf).
- Identified leaf test files each contain the marker helper definition.
- Leaf test files still contain `type Request` (inlined root types).

```go
import (
	"path/filepath"
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
	if resp.MarkerDefCount < 2 {
		t.Fatalf("classic inline should define %s in each leaf (count=%d files=%v gen=%s)",
			experimentP1MarkerFunc, resp.MarkerDefCount, basenames(resp.MarkerDefFiles), resp.GenDir)
	}
	if len(resp.LeafTestFiles) < 2 {
		t.Fatalf("expected ≥2 leaf *_test.go under gen, got %v", basenames(resp.LeafTestFiles))
	}
	for i, leaf := range resp.LeafTestFiles {
		if i < len(resp.LeafHasMarkerDef) && !resp.LeafHasMarkerDef[i] {
			t.Fatalf("classic leaf %s missing %s definition", leaf, experimentP1MarkerFunc)
		}
		if i < len(resp.LeafHasTypeReq) && !resp.LeafHasTypeReq[i] {
			t.Fatalf("classic leaf %s missing inlined type Request", leaf)
		}
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
