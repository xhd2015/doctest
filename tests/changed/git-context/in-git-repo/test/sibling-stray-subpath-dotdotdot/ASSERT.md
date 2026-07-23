## Expected

- Filtered paths are exactly `[leaf_a]`.
- Same policy as sibling-stray-untracked (path argv orthogonal).

```go
import (
	"reflect"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}
	want := []string{"leaf_a"}
	if !reflect.DeepEqual(resp.FilteredPaths, want) {
		t.Fatalf("FilteredPaths = %#v, want %#v", resp.FilteredPaths, want)
	}
}
```
