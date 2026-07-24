---
label: heavy
explanation: nested doctest test quiet prepare on mega parent fixture
---

## Expected

- Exit code 0 (quiet Light→Hydrate prepare succeeds).
- User-visible output reports **1** planned test for the parent tree.
- Output must **not** claim `intermediate/SETUP.md: must have a Go code block`.

## Exit Code

- 0

```go
import (
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = req
	requireExit0(t, resp, err, "quiet-ok")
	out := combinedOutput(resp)
	if hasIntermediateGoBlockError(out) {
		t.Fatalf("quiet must not report intermediate SETUP Go-block error\noutput:\n%s", out)
	}
	n := parsePlannedTests(out)
	if n != 1 {
		t.Fatalf("expected planned 1 parent test under quiet, got %d\noutput:\n%s", n, out)
	}
}
```
