## Expected

- Exit code 0: prose-only SETUP.md (no Go code block) is allowed for intermediate dirs.

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	if err != nil {
		t.Fatalf("run failed unexpectedly: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected zero exit for prose-only SETUP (no go block), got %d\nstdout:\n%s\nstderr:\n%s",
			resp.ExitCode, resp.Stdout, resp.Stderr)
	}
}
```
