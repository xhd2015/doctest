---
explanation: nested doctest CLI + compile of fixture that uses d *session.Doctest
---

## Expected

- Subprocess `doctest test` on the fixture exits 0.
- Fixture leaf reads non-empty `d.DOCTEST_ROOT` / `d.DOCTEST_CASE` / `d.DOCTEST_SESSION_ID`.

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
	if err != nil && resp != nil && resp.ExitCode == 0 {
		t.Fatalf("unexpected run error: %v", err)
	}
	if resp == nil {
		t.Fatal("nil response")
	}
	if resp.ExitCode != 0 {
		t.Fatalf("e2e fixture with d expected PASS, exit=%d\nstdout:\n%s\nstderr:\n%s\nerr=%v",
			resp.ExitCode, resp.Stdout, resp.Stderr, err)
	}
}
```
