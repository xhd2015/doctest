## Expected

- Exit code 0 (all `heavy`, no `e2e` → L3 share 0%).

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
		t.Fatalf("expected zero exit when all leaves are heavy-only (L2), got %d\nstdout:\n%s\nstderr:\n%s",
			resp.ExitCode, resp.Stdout, resp.Stderr)
	}
}
```
