## Expected

- Second write succeeds after parent go.mod change.
- Gen go.mod content differs from first write (includes go 1.20).
- `doctest.tidy-done` is absent (gen go.mod wrote).
- `doctest.gomod-src` still present (rewritten for new fingerprint).

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("WriteGoMod after parent go.mod change failed: %v", err)
	}
	if !strings.Contains(resp.GoModContent, "go 1.20") {
		t.Fatalf("expected regenerated go.mod with go 1.20:\n%s", resp.GoModContent)
	}
	if req.SnapGoModContentBefore != "" && resp.GoModContent == req.SnapGoModContentBefore {
		t.Fatal("gen go.mod content should change after parent go.mod invalidation")
	}
	if resp.TidyDoneExists {
		t.Fatal("expected doctest.tidy-done dropped when gen go.mod rewrote")
	}
	if !resp.GomodSrcExists {
		t.Fatal("expected doctest.gomod-src after rebuild")
	}
}
```
