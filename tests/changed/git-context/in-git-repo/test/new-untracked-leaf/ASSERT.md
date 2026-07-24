## Expected

- Filtered paths are exactly `[leaf_c]`.
- TotalInTree is 3; ChangedCount is 1.

```go
import (
	"reflect"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}
	want := []string{"leaf_c"}
	if !reflect.DeepEqual(resp.FilteredPaths, want) {
		t.Fatalf("FilteredPaths = %#v, want %#v", resp.FilteredPaths, want)
	}
	if resp.Info.TotalInTree != 3 {
		t.Fatalf("TotalInTree = %d, want 3", resp.Info.TotalInTree)
	}
	if resp.Info.ChangedCount != 1 {
		t.Fatalf("ChangedCount = %d, want 1", resp.Info.ChangedCount)
	}
}
```
