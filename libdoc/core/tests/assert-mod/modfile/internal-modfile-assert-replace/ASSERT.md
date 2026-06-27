## Expected

- Modfile contains parent `go.mod` contents (including existing replace).
- Modfile appends `replace github.com/xhd2015/doctest/assert => <cache>`.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("WriteInternalModfile failed: %v", err)
	}
	modfile := resp.ModfileContent
	if !strings.Contains(modfile, "module example.com/app") {
		t.Fatalf("expected parent module declaration preserved, got:\n%s", modfile)
	}
	if !strings.Contains(modfile, "replace example.com/dep => ../dep") {
		t.Fatalf("expected parent replace preserved, got:\n%s", modfile)
	}
	needle := "replace github.com/xhd2015/doctest/assert =>"
	if !strings.Contains(modfile, needle) {
		t.Fatalf("expected assert replace appended, got:\n%s", modfile)
	}
	if !strings.Contains(modfile, req.AssertCacheDir) {
		t.Fatalf("expected assert replace to point at %s, got:\n%s", req.AssertCacheDir, modfile)
	}
}
```