## Expected

- `resp.Display` is `"./child"` (platform-native separators).

## Errors

- `err` is nil.

```go
import (
	"path/filepath"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	want := "." + string(filepath.Separator) + "child"
	if resp.Display != want {
		t.Fatalf("expected %q, got %q (cwd=%q path=%q)", want, resp.Display, resp.Cwd, req.Path)
	}
}
```