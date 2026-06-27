## Expected

- Generated nested `go.mod` contains `module testcase`.
- Contains `replace github.com/xhd2015/doctest/assert =>` pointing at cache dir.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("WriteGoMod failed: %v", err)
	}
	goMod := resp.GoModContent
	if !strings.Contains(goMod, "module testcase") {
		t.Fatalf("expected module testcase, got:\n%s", goMod)
	}
	needle := "replace github.com/xhd2015/doctest/assert =>"
	if !strings.Contains(goMod, needle) {
		t.Fatalf("expected assert replace, got:\n%s", goMod)
	}
	if !strings.Contains(goMod, req.AssertCacheDir) {
		t.Fatalf("expected replace to point at %s, got:\n%s", req.AssertCacheDir, goMod)
	}
}
```