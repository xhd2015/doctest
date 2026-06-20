## Expected

- `resp.Display` is `"./child"`.

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
		t.Fatalf("expected %q, got %q for relative input %q", want, resp.Display, req.Path)
	}
}
```