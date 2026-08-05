## Expected

- Exit code 0 (10 leaves, 1 e2e → 10% L3 share is within MaxL3Pct=10).

## Exit Code

- 0

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	if err != nil {
		t.Fatalf("run failed unexpectedly: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected zero exit within L3 budget (10 leaves, 1 e2e), got %d\nstdout:\n%s\nstderr:\n%s",
			resp.ExitCode, resp.Stdout, resp.Stderr)
	}
}
```
