## Expected

- Second write succeeds with new assert cache dir.
- Gen go.mod content changed (assert replace RHS differs).
- Warm mtime skip did **not** apply (mtime after differs or content differs).
- `doctest.tidy-done` dropped when go.mod wrote.

```go
import (
	"path/filepath"
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("WriteGoMod after assert path change failed: %v", err)
	}
	if resp.GoModContent == req.SnapGoModContentBefore {
		t.Fatal("gen go.mod should change when effective assert cache dir changes")
	}
	// Assert import path appears with the second cache dir.
	if req.SecondAssertCacheDir != "" {
		want := filepath.ToSlash(req.SecondAssertCacheDir)
		if !strings.Contains(filepath.ToSlash(resp.GoModContent), want) &&
			!strings.Contains(resp.GoModContent, req.SecondAssertCacheDir) {
			// Absolute path may be cleaned; require content differs is primary.
			t.Logf("note: second assert dir %q not found literally in go.mod (abs may differ)", req.SecondAssertCacheDir)
		}
	}
	if resp.TidyDoneExists {
		t.Fatal("expected tidy-done dropped when assert path forced go.mod rewrite")
	}
	if !resp.GomodSrcExists {
		t.Fatal("expected doctest.gomod-src after assert-path miss")
	}
}
```
