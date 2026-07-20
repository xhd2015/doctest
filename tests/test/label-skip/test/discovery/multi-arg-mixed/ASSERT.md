---
label: heavy
---

## Expected

- Exit 0, PASS(2/2), aggregated skip block lists labeled leaves from discovery pass.

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
	assertSkipCompact(t, resp.Stdout, 2, map[string]int{"ui-automation": 2})
	assertResultSummary(t, resp.Stdout, 2, 2)
}
```