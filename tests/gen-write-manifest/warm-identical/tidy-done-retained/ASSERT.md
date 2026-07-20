## Expected

- Second WriteGoMod succeeds.
- `doctest.tidy-done` still exists (module files did not actually write).
- Unified layout: manifest present, no `doctest.gomod-fp`.

```go
import (
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("second WriteGoMod failed: %v", err)
	}
	requireUnifiedLayout(t, resp)
	if !resp.TidyDoneExists {
		t.Fatal("expected doctest.tidy-done retained when go.mod/go.sum did not write")
	}
}
```
