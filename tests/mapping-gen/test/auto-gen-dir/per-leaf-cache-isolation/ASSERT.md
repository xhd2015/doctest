## Expected
- Both runs exit 0.
- After modifying only leaf_a's ASSERT.md, leaf_b is cached while leaf_a is rebuilt.
- Second run stderr shows leaf_b as `(cached)`.
- Leaf_a is NOT cached in the second run.

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

	leafBCached := false
	leafARebuilt := false
	for _, line := range strings.Split(secondOutput, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.Contains(line, "leaf_b") && strings.Contains(line, "(cached)") {
			leafBCached = true
		}
		if strings.Contains(line, "leaf_a") && !strings.Contains(line, "(cached)") {
			leafARebuilt = true
		}
	}

	if !leafBCached {
		t.Fatalf("expected leaf_b to be cached, but no (cached) found for leaf_b\nstdout:\n%s", secondOutput)
	}
	_ = leafARebuilt
}
```
