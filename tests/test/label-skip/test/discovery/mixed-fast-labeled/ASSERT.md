---
label: heavy
explanation: CLI skip/summary contract via doctest binary
---

## Expected

- Exit 0, one test runs, labeled leaf skipped with exact skip block and PASS(1/1).

```go
import (
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("exit code = %d\nstdout:\n%s\nstderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	tree := req.Args[1]
	assertSkipBlockExact(t, resp.Stdout, tree, "labeled_leaf", "ui-automation", "heavy ui test")
	assertResultSummary(t, resp.Stdout, 1, 1)
}
```