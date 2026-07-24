## Expected

- Same selection as assert-only: `[leaf_a]`, detail `1 leaf`.
- Documents that `./...` path discovery does not change filter policy.

```go
import (
	"reflect"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}
	want := []string{"leaf_a"}
	if !reflect.DeepEqual(resp.FilteredPaths, want) {
		t.Fatalf("FilteredPaths = %#v, want %#v", resp.FilteredPaths, want)
	}
	if resp.Info.ChangedCount != 1 || resp.Info.Detail != "1 leaf" {
		t.Fatalf("info = %#v, want ChangedCount=1 Detail=%q", resp.Info, "1 leaf")
	}
}
```
