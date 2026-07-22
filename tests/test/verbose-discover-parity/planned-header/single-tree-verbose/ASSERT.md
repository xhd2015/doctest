---
label: heavy
explanation: nested doctest test -v single-tree planned header lock
---

## Expected

- Exit code 0.
- Combined stdout/stderr contains a **planned tests** signal for **1** test
  before/alongside the go test invocation — e.g. `doctest: … (1 tests)`,
  `─── 1 test cases`, or `1 tests planned`.
- Stderr includes `go test` with `-v` (verbose still runs go test -v).

## Exit Code

- 0

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	_ = req
	requireExit0(t, resp, err, "single-tree-verbose")
	out := combinedOutput(resp)
	n := parsePlannedTests(out)
	if n != 1 {
		t.Fatalf("under -v, expected planned 1 test in user-visible output, got %d\noutput:\n%s", n, out)
	}
	// Prefer evidence the count is associated with doctest announce / planned, not only PASS (1/1).
	// parsePlannedTests already excludes bare PASS (1/1) by requiring "tests" / "test cases" / "planned".
	if !strings.Contains(resp.Stderr, "go test") {
		t.Fatalf("expected stderr to show go test invocation under -v:\n%s", resp.Stderr)
	}
}
```
