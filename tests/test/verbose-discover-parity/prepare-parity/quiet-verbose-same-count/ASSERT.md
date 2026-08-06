---
explanation: dual nested doctest test quiet vs -v case-count parity
---

## Expected

- Quiet subprocess exit 0; verbose subprocess exit 0.
- Planned test count from quiet output equals planned count from verbose output.
- Both report **1** (parent `own_leaf` only; nested tree is not part of parent target).
- Verbose output has no intermediate SETUP Go-block hard-fail.

## Exit Code

- 0 for both modes

```go
import (
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = resp
	if err != nil {
		t.Fatalf("dual_cli harness error: %v", err)
	}
	if req.Quiet == nil || req.Verbose == nil {
		t.Fatalf("dual_cli must fill req.Quiet and req.Verbose (quiet=%v verbose=%v)",
			req.Quiet != nil, req.Verbose != nil)
	}
	requireExit0(t, req.Quiet, nil, "quiet half")
	requireExit0(t, req.Verbose, nil, "verbose half")

	qOut := combinedOutput(req.Quiet)
	vOut := combinedOutput(req.Verbose)
	if hasIntermediateGoBlockError(vOut) {
		t.Fatalf("verbose half must not fail intermediate SETUP Go-block\nverbose:\n%s", vOut)
	}
	pq := parsePlannedTests(qOut)
	pv := parsePlannedTests(vOut)
	if pq != 1 {
		t.Fatalf("quiet planned want 1, got %d\nquiet:\n%s", pq, qOut)
	}
	if pv != 1 {
		t.Fatalf("verbose planned want 1, got %d\nverbose:\n%s", pv, vOut)
	}
	if pq != pv {
		t.Fatalf("quiet planned %d != verbose planned %d\nquiet:\n%s\nverbose:\n%s",
			pq, pv, qOut, vOut)
	}
}
```
