## Expected

- Second WriteGoMod succeeds after modules.txt-only change.
- Gen go.mod includes the new vendored module.
- `doctest.tidy-done` is absent (gen go.mod actually wrote).
- Unified layout + gomod-src present after rewrite.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("WriteGoMod after modules.txt change failed: %v", err)
	}
	requireUnifiedLayout(t, resp)
	if !strings.Contains(resp.GoModContent, "example.com/other") {
		t.Fatalf("expected regenerated go.mod to include example.com/other:\n%s", resp.GoModContent)
	}
	if !strings.Contains(req.SnapGoModContentBefore, "example.com/dep") {
		t.Fatalf("precondition: first go.mod should have included dep:\n%s", req.SnapGoModContentBefore)
	}
	if resp.TidyDoneExists {
		t.Fatal("expected doctest.tidy-done removed when modules.txt forced go.mod rewrite")
	}
	if !resp.GomodSrcExists {
		t.Fatal("expected doctest.gomod-src rewritten after miss")
	}
}
```
