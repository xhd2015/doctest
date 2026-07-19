---
label: heavy
explanation: nested doctest CLI; fixture reads leaf-local file via d.DOCTEST_CASE
---

## Expected

- Subprocess exits 0.
- Fixture successfully read `fixture.txt` through `d.DOCTEST_CASE` (not bare relative cwd).

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
		t.Fatalf("leaf-local via d.DOCTEST_CASE expected PASS, exit=%d\nstdout:\n%s\nstderr:\n%s\nerr=%v",
			resp.ExitCode, resp.Stdout, resp.Stderr, err)
	}
}
```
