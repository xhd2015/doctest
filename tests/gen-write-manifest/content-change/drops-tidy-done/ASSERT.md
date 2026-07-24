## Expected

- Second WriteGoMod succeeds.
- `doctest.tidy-done` is **absent** after go.mod actually wrote.
- Unified layout still holds (manifest present, no gomod-fp).

```go
import (
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("WriteGoMod after source change failed: %v", err)
	}
	requireUnifiedLayout(t, resp)
	if resp.TidyDoneExists {
		t.Fatal("expected doctest.tidy-done removed when go.mod actually wrote")
	}
}
```
