## Expected

- No new `*.jsonl` files under MetricsRoot after the suite.
- Suite failure is not required; if RunErr is set solely for unrelated reasons, still require empty RunFiles.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected harness error: %v", err)
	}
	if len(resp.RunFiles) != 0 {
		t.Fatalf("default MetricsOn=false should write no run files; got %v\nstderr:\n%s", resp.RunFiles, resp.Stderr)
	}
}
```
