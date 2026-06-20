## Expected

- `resp.Display` starts with `"~/"` or `"~\\"` (platform-native).
- `resp.Display` contains `"mapping-gen"`.
- `resp.Display` does **not** contain the raw absolute home prefix.

## Errors

- `err` is nil.

```go
import (
	"os"
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(resp.Display, "~") {
		t.Fatalf("expected ~ prefix, got %q", resp.Display)
	}
	if strings.Contains(resp.Display, home) {
		t.Fatalf("display must not contain raw home %q: %q", home, resp.Display)
	}
	if !strings.Contains(resp.Display, "mapping-gen") {
		t.Fatalf("expected mapping-gen in display path, got %q", resp.Display)
	}
}
```