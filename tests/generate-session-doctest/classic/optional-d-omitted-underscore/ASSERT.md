## Expected

- Assemble fails when author omits `d *session.Doctest`.
- Error message mentions missing `d` and/or no auto-inject (not a successful gen with `_`).

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err == nil {
		t.Fatal("expected assemble error when author omits d *session.Doctest; got success (no auto-inject)")
	}
	msg := err.Error()
	if !strings.Contains(msg, "*session.Doctest") && !strings.Contains(msg, "no auto-inject") && !strings.Contains(msg, "missing d") {
		t.Fatalf("expected clear missing-d error, got: %v", err)
	}
}
```
