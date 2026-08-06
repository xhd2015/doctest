## Expected

- Validation succeeds.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("Run err: %v", err)
	}
	if resp.ErrMsg != "" {
		t.Fatalf("want in-range OK, got %q", resp.ErrMsg)
	}
	if resp.Min != "1.0" || resp.Max != "99.0" {
		t.Fatalf("bounds min=%q max=%q", resp.Min, resp.Max)
	}
}
```
