## Expected

- `resp.Display` is `"./a/b/c"` (platform-native separators).

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
	want := "." + string(filepath.Separator) + "a" + string(filepath.Separator) + "b" + string(filepath.Separator) + "c"
	if resp.Display != want {
		t.Fatalf("expected %q, got %q", want, resp.Display)
	}
}
```