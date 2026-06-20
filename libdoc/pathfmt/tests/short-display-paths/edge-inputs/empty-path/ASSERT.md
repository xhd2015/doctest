## Expected

- `resp.Display` is `""` (unchanged).

## Errors

- `err` is nil.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if resp.Display != "" {
		t.Fatalf("expected empty display for empty input, got %q", resp.Display)
	}
}
```