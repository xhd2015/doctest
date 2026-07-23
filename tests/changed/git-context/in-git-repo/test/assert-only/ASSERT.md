## Expected

- Filtered paths are exactly `[leaf_a]`.
- `ChangedCount` is 1; detail is `1 leaf`.
- Announce is true (non-zero changed).

## Exit Code

n/a (L2 policy)

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
	if resp.Info.TotalInTree != 2 {
		t.Fatalf("TotalInTree = %d, want 2", resp.Info.TotalInTree)
	}
	if resp.Info.ChangedCount != 1 {
		t.Fatalf("ChangedCount = %d, want 1", resp.Info.ChangedCount)
	}
	if resp.Info.Detail != "1 leaf" {
		t.Fatalf("Detail = %q, want %q", resp.Info.Detail, "1 leaf")
	}
	if !resp.Announce {
		t.Fatal("expected Announce=true for non-zero changed")
	}
}
```
