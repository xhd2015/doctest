## Expected

- Generated source assigns `DOCTEST_CASE` (on `d`) to `filepath.Join(DocTestRoot, CasePath)`.
- Inject contract still holds.

```go
import (
	"path/filepath"
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("classic assemble failed: %v", err)
	}
	assertInjectContract(t, "classic-case-path", resp.Source)
	wantCase := filepath.Join(req.DocTestRoot, req.CasePath)
	if !containsCasePathAssignment(resp.Source, wantCase) {
		t.Fatalf("expected DOCTEST_CASE to embed abs leaf path %q\n%s", wantCase, resp.Source)
	}
	// ROOT should also appear as the doc root.
	if !strings.Contains(resp.Source, req.DocTestRoot) {
		t.Fatalf("expected DOCTEST_ROOT abs path %q in source\n%s", req.DocTestRoot, resp.Source)
	}
}
```
