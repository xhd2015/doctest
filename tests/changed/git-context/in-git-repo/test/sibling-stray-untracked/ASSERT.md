## Expected

- Filtered paths are exactly `[leaf_a]`.
- Sibling stray.go does not select `leaf_b`.

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
		t.Fatalf("FilteredPaths = %#v, want %#v (sibling stray must not widen)", resp.FilteredPaths, want)
	}
	if resp.Info.ChangedCount != 1 {
		t.Fatalf("ChangedCount = %d, want 1", resp.Info.ChangedCount)
	}
}
```
