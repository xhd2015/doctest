## Expected

- Suite run succeeds.
- Exactly one generated `*.go` file defines `func ExperimentP1RootMarker`
  (root package owns the helper once; leaves do not redefine it).

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
		t.Fatalf("hierarchical RunTest failed before layout check: %s\nstderr:\n%s",
			resp.RunErr, resp.Stderr)
	}
	if resp.MarkerDefCount != 1 {
		names := make([]string, len(resp.MarkerDefFiles))
		for i, p := range resp.MarkerDefFiles {
			names[i] = filepath.ToSlash(p)
		}
		t.Fatalf("want exactly 1 definition of %s under gen, got count=%d files=%v gen=%s",
			experimentP1MarkerFunc, resp.MarkerDefCount, names, resp.GenDir)
	}
}
```
