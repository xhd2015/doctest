## Expected

- `resp.Display` equals the absolute path (unchanged aside from normalization).
- `resp.Display` does **not** start with `"~"`.
- `resp.Display` does **not** start with `"."` (not cwd-relative).

## Errors

- `err` is nil.

```go
import (
	"path/filepath"
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	abs, err := filepath.Abs(req.Path)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Display != abs {
		t.Fatalf("expected absolute %q, got %q", abs, resp.Display)
	}
	if strings.HasPrefix(resp.Display, "~") {
		t.Fatalf("outside-home path must not use ~ prefix: %q", resp.Display)
	}
	if strings.HasPrefix(resp.Display, ".") {
		t.Fatalf("outside-home path must not be cwd-relative: %q", resp.Display)
	}
}
```