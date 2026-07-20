## Expected

- Second WriteGoMod succeeds.
- `go.mod` mtime unchanged under ineffective assert/session flag differences.
- `doctest.tidy-done` retained.
- Generated go.mod has parent replace but **no** assert/session submodule replaces.
- Unified layout: `doctest.gen-manifest` present, `doctest.gomod-fp` absent.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("second WriteGoMod failed: %v", err)
	}
	requireUnifiedLayout(t, resp)
	if resp.GoModMtimeBefore.IsZero() {
		t.Fatal("missing go.mod mtime snapshot")
	}
	if !resp.GoModMtimeBefore.Equal(resp.GoModMtimeAfter) {
		t.Fatalf("go.mod mtime changed under ineffective assert/session flags: %v -> %v",
			resp.GoModMtimeBefore, resp.GoModMtimeAfter)
	}
	if !resp.TidyDoneExists {
		t.Fatal("expected tidy-done retained when ineffective flags do not change go.mod")
	}
	if strings.Contains(resp.GoModContent, "replace github.com/xhd2015/doctest/assert =>") {
		t.Fatalf("assert replace must not appear for doctest self-module:\n%s", resp.GoModContent)
	}
	if strings.Contains(resp.GoModContent, "replace github.com/xhd2015/doctest/session =>") {
		t.Fatalf("session replace must not appear for doctest self-module:\n%s", resp.GoModContent)
	}
	if !strings.Contains(resp.GoModContent, "replace github.com/xhd2015/doctest =>") {
		t.Fatalf("expected parent module replace, got:\n%s", resp.GoModContent)
	}
}
```
