## Expected

- Filtered paths empty.
- `ChangedCount` 0; `Announce` false (silent without verbose).

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}
	if len(resp.FilteredPaths) != 0 {
		t.Fatalf("FilteredPaths = %#v, want empty", resp.FilteredPaths)
	}
	if resp.Info.TotalInTree != 2 {
		t.Fatalf("TotalInTree = %d, want 2", resp.Info.TotalInTree)
	}
	if resp.Info.ChangedCount != 0 {
		t.Fatalf("ChangedCount = %d, want 0", resp.Info.ChangedCount)
	}
	if resp.Announce {
		t.Fatal("Announce should be false for zero-change without verbose")
	}
}
```
