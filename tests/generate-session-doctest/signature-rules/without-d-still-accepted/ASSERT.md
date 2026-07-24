## Expected

- Legacy without-d signatures are rejected (`ParseErr` non-empty).

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected run error: %v", err)
	}
	if resp == nil {
		t.Fatal("nil response")
	}
	if resp.ParseErr == "" {
		t.Fatal("expected without-d signatures rejected; ParseErr empty")
	}
	if !strings.Contains(resp.ParseErr, "no auto-inject") && !strings.Contains(resp.ParseErr, "d *session.Doctest") {
		t.Fatalf("expected rejection message about required d, got: %s", resp.ParseErr)
	}
}
```
