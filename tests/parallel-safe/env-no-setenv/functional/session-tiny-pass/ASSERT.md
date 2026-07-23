---
label: e2e, heavy
explanation: nested doctest test generate + go test for session inject smoke
---

## Expected

- Subprocess `doctest test` exits 0.
- Nested leaf Assert passed (non-empty `d.DOCTEST_SESSION_ID` from suite child env).
- GREEN before and after P1 product fix (implementer keeps session via `cmd.Env`).

## Exit Code

- 0

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	if err != nil {
		t.Fatalf("run failed: %v\nstderr:\n%s", err, resp.Stderr)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected exit 0 for session tiny fixture, got %d\nstdout:\n%s\nstderr:\n%s",
			resp.ExitCode, resp.Stdout, resp.Stderr)
	}
}
```
