## Expected

- ShouldWarn is false for non-default (label-all) suite.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.ShouldWarn {
		t.Fatal("ShouldWarn = true, want false for label-all / non-default suite")
	}
}
```
