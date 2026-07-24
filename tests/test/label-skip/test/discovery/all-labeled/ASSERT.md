---
label: heavy
explanation: CLI skip/summary contract via doctest binary
---

## Expected

- Exit 0, exact skip block, no PASS/FAIL summary line.

```go
import (
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("exit code = %d\nstdout:\n%s\nstderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	tree := req.Args[1]
	assertSkipBlockExact(t, resp.Stdout, tree, "labeled_leaf", "human-guided-ui-test", "manual only")
	assertNoResultSummary(t, resp.Stdout)
}
```