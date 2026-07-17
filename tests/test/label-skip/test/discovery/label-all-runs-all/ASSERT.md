---
label: heavy
---

## Expected

- Exit 0, both leaves run, no skip block, PASS(2/2).

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
	assertNoSkipBlock(t, resp.Stdout)
	assertResultSummary(t, resp.Stdout, 2, 2)
}
```
