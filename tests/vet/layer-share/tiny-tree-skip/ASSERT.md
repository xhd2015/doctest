## Expected

- Exit code 0 (share enforcement only when leaves ≥ 10).

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
		t.Fatalf("expected zero exit for tiny tree (3 leaves, 1 e2e; MinLeaves=10), got %d\nstdout:\n%s\nstderr:\n%s",
			resp.ExitCode, resp.Stdout, resp.Stderr)
	}
}
```
