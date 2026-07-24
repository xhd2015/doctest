## Expected

- ShouldWarn is false when the suite is label-filtered (not default).

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.ShouldWarn {
		t.Fatal("ShouldWarn = true, want false for label-filtered suite")
	}
}
```
