## Expected

- Both warn cases return false.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.WarnResults) != len(req.WarnCases) {
		t.Fatalf("results=%d cases=%d", len(resp.WarnResults), len(req.WarnCases))
	}
	for i, c := range req.WarnCases {
		if resp.WarnResults[i] != c.Want {
			t.Fatalf("case %q: ShouldWarn=%v, want %v (elapsed=%v)", c.Name, resp.WarnResults[i], c.Want, c.Elapsed)
		}
	}
}
```
