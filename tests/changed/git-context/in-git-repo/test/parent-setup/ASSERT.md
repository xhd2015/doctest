## Expected

- Filtered paths are both shared leaves (order sorted).
- `ChangedCount` is 2.
- Detail mentions `shared/SETUP.md` affecting other tests.

```go
import (
	"reflect"
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}
	want := []string{"shared/leaf_a", "shared/leaf_b"}
	if !reflect.DeepEqual(resp.FilteredPaths, want) {
		t.Fatalf("FilteredPaths = %#v, want %#v", resp.FilteredPaths, want)
	}
	if resp.Info.ChangedCount != 2 {
		t.Fatalf("ChangedCount = %d, want 2", resp.Info.ChangedCount)
	}
	if !strings.Contains(resp.Info.Detail, "shared/SETUP.md") {
		t.Fatalf("Detail = %q, want shared/SETUP.md mention", resp.Info.Detail)
	}
}
```
