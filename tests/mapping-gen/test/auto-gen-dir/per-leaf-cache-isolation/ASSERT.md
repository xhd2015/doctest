---
label: heavy
---

## Expected
- Both runs exit 0.
- After modifying only leaf_a's ASSERT.md, leaf_b is cached while leaf_a is rebuilt.
- Second run stdout summary reports cached packages (`N Cached`).

## Exit Code
- Exit code 0.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("run failed: %v\nstderr:\n%s", err, resp.Stderr)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected exit 0, got %d\nstderr:\n%s", resp.ExitCode, resp.Stderr)
	}

	if firstResp == nil || secondResp == nil {
		t.Fatal("multi-run state not set by Run")
	}
	if secondResp.ExitCode != 0 {
		t.Fatalf("second run exit %d\nstderr:\n%s", secondResp.ExitCode, secondResp.Stderr)
	}

	secondOutput := secondResp.Stdout
	if !strings.Contains(secondOutput, "Cached") {
		t.Fatalf("expected second run summary to report cached packages, stdout:\n%s", secondOutput)
	}
	// Unified suite: may report "1 Pass, 1 Cached" (leaf counts) rather than "2 Pass".
	if !strings.Contains(secondOutput, "2 Pass") &&
		!(strings.Contains(secondOutput, "Pass") && strings.Contains(secondOutput, "Cached")) &&
		!strings.Contains(secondOutput, "PASS (2/2)") {
		t.Fatalf("expected both leaves to pass on second run (2 Pass or Pass+Cached), stdout:\n%s", secondOutput)
	}
}
```
