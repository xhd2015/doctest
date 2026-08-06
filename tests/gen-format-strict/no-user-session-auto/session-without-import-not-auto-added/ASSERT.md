## Expected

- `WriteFormattedGo` does **not** inject `github.com/xhd2015/doctest/session`.
- Written source still contains `session.` usage without a session import.
- `go build` on the written package fails (`BuildErr` non-empty).
- Documents: no user-path auto-inject of session via format reconcile maps.

```go
import (
	"strings"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	if err != nil {
		t.Fatalf("harness error: %v", err)
	}
	if resp == nil {
		t.Fatal("nil response")
	}
	src := resp.FormattedSource
	if src == "" {
		t.Fatal("A5: empty FormattedSource")
	}
	if containsSessionImport(src) {
		t.Fatalf("A5: WriteFormattedGo auto-added session import for user session. usage:\n%s", src)
	}
	if !strings.Contains(src, "session.") {
		t.Fatalf("A5: expected session. selector preserved in source:\n%s", src)
	}
	if resp.BuildErr == "" {
		t.Fatalf("A5: expected go build to fail without session import, BuildErr empty\nout:\n%s\nsrc:\n%s",
			resp.BuildOutput, src)
	}
}
```
