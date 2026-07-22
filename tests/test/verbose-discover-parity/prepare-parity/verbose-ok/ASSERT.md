---
label: heavy
explanation: nested doctest test -v prepare on mega parent (RED until parity fix)
---

## Expected

- Exit code 0 under `-v` (prepare must not abort the parent tree).
- Combined output must **not** contain `intermediate` + `SETUP.md` +
  `must have a Go code block` (or equivalent hard prepare fail).
- Planned test count for the parent is **1**.

## Errors

- Today (before implement): prepare fails with full-discover validation on the
  intermediate path even though Light already selected `own_leaf`.

## Exit Code

- 0

```go
import (
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	_ = req
	requireExit0(t, resp, err, "verbose-ok")
	out := combinedOutput(resp)
	if hasIntermediateGoBlockError(out) {
		t.Fatalf("-v prepare must not fail with intermediate SETUP Go-block error\noutput:\n%s", out)
	}
	n := parsePlannedTests(out)
	if n != 1 {
		t.Fatalf("expected planned 1 parent test under -v, got %d\noutput:\n%s", n, out)
	}
}
```
