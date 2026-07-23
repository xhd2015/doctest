## Expected

- Markdown path list is exactly `leaf_a/ASSERT.md`.
- Root `DOCTEST.md` is not selected (would have failed Version check if validated).

```go
import (
	"reflect"
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}
	want := []string{"leaf_a/ASSERT.md"}
	if !reflect.DeepEqual(resp.MarkdownPaths, want) {
		t.Fatalf("MarkdownPaths = %#v, want %#v", resp.MarkdownPaths, want)
	}
	for _, p := range resp.MarkdownPaths {
		if strings.Contains(p, "DOCTEST.md") {
			t.Fatalf("root DOCTEST.md must not be selected: %#v", resp.MarkdownPaths)
		}
	}
}
```
