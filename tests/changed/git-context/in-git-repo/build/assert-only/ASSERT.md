## Expected

- Filtered paths are exactly `[leaf_a]` — build shares test selection API.
- Detail is `1 leaf`.

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
		t.Fatalf("info = %#v", resp.Info)
	}
}
```
