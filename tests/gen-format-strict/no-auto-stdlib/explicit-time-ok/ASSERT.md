## Expected

- Generate + suite `go test` succeed (`resp.RunErr` empty).
- Generated leaf source may mention `time` / `"time"` (author imported it).

```go
import (
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
	if resp.RunErr != "" {
		t.Fatalf("expected success with explicit import \"time\", got RunErr=%v\nstdout:\n%s\nstderr:\n%s\nleaf:\n%s",
			resp.RunErr, resp.Stdout, resp.Stderr, resp.LeafGo)
	}
}
```
