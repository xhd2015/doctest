## Expected

- Second write succeeds.
- `doctest.tidy-done` still exists after warm hit.

```go
import (
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("warm WriteGoModWithVendorBridges failed: %v", err)
	}
	if !resp.TidyDoneExists {
		t.Fatal("expected doctest.tidy-done retained on warm hit")
	}
}
```
