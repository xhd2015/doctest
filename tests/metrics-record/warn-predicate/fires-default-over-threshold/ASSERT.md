## Expected

- `ShouldWarnDefaultSuiteSlow` returns true.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resp.ShouldWarn {
		t.Fatalf("ShouldWarn = false, want true (default suite, total=%d, elapsed=%v)", req.Total, req.Elapsed)
	}
}
```
