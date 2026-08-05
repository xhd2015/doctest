## Expected

- Non-zero exit.
- stderr signals Error / missing path (fatal discovery).

## Exit Code

- non-zero

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	requireFail(t, resp, err)
	se := strings.ToLower(resp.Stderr)
	ok := strings.Contains(se, "error") ||
		strings.Contains(se, "no such file") ||
		strings.Contains(se, "not exist") ||
		strings.Contains(se, "does-not-exist-list-root") ||
		strings.Contains(se, "stat")
	if !ok {
		t.Fatalf("expected missing-path error on stderr:\n%s", resp.Stderr)
	}
}
```
