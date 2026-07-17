## Expected

- ShouldWarn is false when total is 0.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.ShouldWarn {
		t.Fatal("ShouldWarn = true, want false when total=0")
	}
}
```
