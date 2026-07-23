## Expected

- Filtered paths empty despite a changed file outside the tree.
- `ChangedCount` 0; `Announce` false.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}
	if len(resp.FilteredPaths) != 0 {
		t.Fatalf("FilteredPaths = %#v, want empty", resp.FilteredPaths)
	}
	if resp.Info.ChangedCount != 0 {
		t.Fatalf("ChangedCount = %d, want 0", resp.Info.ChangedCount)
	}
	if resp.Announce {
		t.Fatal("Announce should be false when no doctest leaves match")
	}
}
```
