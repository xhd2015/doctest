## Expected

- Exit code 0.
- Output refers to the newer run (stem or run_id containing `newrun01` or full newer basename stem).
- Output reflects latest suite stats (e.g. total/passed 4, or slow-leaf path), not solely the older-only leaf.

## Exit Code

- 0

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("exit=%d stderr=%s stdout=%s", resp.ExitCode, resp.Stderr, resp.Stdout)
	}
	out := combinedOut(resp)
	if !strings.Contains(out, "newrun01") && !strings.Contains(out, "2026-07-16-09-00-00-00-newrun01") {
		t.Fatalf("last should identify newest run (newrun01):\n%s", out)
	}
	// Should not be exclusively the old run marker without the new one.
	if strings.Contains(out, "old/only-leaf") && !strings.Contains(out, "slow-leaf") && !strings.Contains(out, "newrun01") {
		t.Fatalf("last appears to summarize older run only:\n%s", out)
	}
}
```
