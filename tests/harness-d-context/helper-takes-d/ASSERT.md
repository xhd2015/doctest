---
explanation: nested doctest CLI; fixture package helper takes d *session.Doctest
---

## Expected

- Subprocess exits 0.
- Fixture helper `joinCase(d, name)` produced the absolute leaf-local path.

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
		t.Fatalf("helper-takes-d expected PASS, exit=%d\nstdout:\n%s\nstderr:\n%s\nerr=%v",
			resp.ExitCode, resp.Stdout, resp.Stderr, err)
	}
}
```
