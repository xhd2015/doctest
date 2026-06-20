## Expected

- `resp.Display` is **not** a cwd-relative `./..` form.
- `resp.Display` starts with `"~"` (parent is under home) **or** equals the absolute parent path.

## Errors

- `err` is nil.

```go
import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if strings.HasPrefix(resp.Display, "."+string(filepath.Separator)+"..") {
		t.Fatalf("parent path must not display as %q, got %q", "./..", resp.Display)
	}
	absParent, _ := filepath.Abs(req.Path)
	home, _ := os.UserHomeDir()
	homeShort := "~" + strings.TrimPrefix(absParent, home)
	if resp.Display != absParent && resp.Display != homeShort {
		t.Fatalf("expected %q or %q, got %q", absParent, homeShort, resp.Display)
	}
}
```