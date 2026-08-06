---
explanation: nested workspace doctest test -v planned trees+tests header
---

## Expected

- Exit code 0.
- User-visible output contains a workspace planned line:
  `doctest: workspace (2 trees, 2 tests)` (or `workspace hub`).
- Planned line appears **before** the `cd … && go test` command line when both
  are present (verbose still prints the cd/go test line).

## Exit Code

- 0

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = req
	requireExit0(t, resp, err, "workspace-verbose")
	out := combinedOutput(resp)
	trees, tests := parseWorkspacePlanned(out)
	if trees != 2 || tests != 2 {
		t.Fatalf("under -v workspace, expected (2 trees, 2 tests), got trees=%d tests=%d\noutput:\n%s",
			trees, tests, out)
	}
	// Ordering: planned workspace line before "cd " go-test announce when both exist.
	planIdx := plannedWorkspaceRe.FindStringIndex(out)
	cdIdx := strings.Index(out, "cd ")
	if planIdx != nil && cdIdx >= 0 && planIdx[0] > cdIdx {
		t.Fatalf("planned workspace line must appear before cd/go test under -v\noutput:\n%s", out)
	}
}
```
