## Expected

- Markdown path list is exactly `leaf_b/SETUP.md`.
- Unchanged sibling `leaf_a` is not selected for vet.

```go
import (
	"reflect"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}
	want := []string{"leaf_b/SETUP.md"}
	if !reflect.DeepEqual(resp.MarkdownPaths, want) {
		t.Fatalf("MarkdownPaths = %#v, want %#v", resp.MarkdownPaths, want)
	}
}
```
