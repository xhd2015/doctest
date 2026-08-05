## Expected
- The command succeeds with exit code 0 (read-only stdio use is not an anti-pattern).

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("run failed unexpectedly: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected zero exit for Fprint(os.Stdout) without reassignment, got %d\nstdout:\n%s\nstderr:\n%s",
			resp.ExitCode, resp.Stdout, resp.Stderr)
	}
}
```
