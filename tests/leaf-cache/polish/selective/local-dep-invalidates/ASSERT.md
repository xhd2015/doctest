---
label: heavy
explanation: nested 3× doctest test; local package mutation mid-sequence
---

## Expected

- Run1 and run2 exit 0; run2 has **Cached > 0**.
- After local package mutation, run3 has **0 Cached** (leaf re-executed).
- Run3 may fail (assert still expects Answer==42 while helper returns 99) — that
  failure **proves** re-execution; do not require exit 0 on run3.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("run1 exit %d\n%s\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	if resp.ExitCode2 != 0 {
		t.Fatalf("run2 exit %d\n%s\n%s", resp.ExitCode2, resp.Stdout2, resp.Stderr2)
	}
	if !stdoutCachedPositive(resp.Stdout2) {
		t.Fatalf("run2 warm expected Cached > 0; got %d\n%s", cachedCount(resp.Stdout2), resp.Stdout2)
	}
	if !stdoutCachedZero(resp.Stdout3) {
		t.Fatalf("after local dep edit expected 0 Cached (re-run); got %d\nstdout3:\n%s", cachedCount(resp.Stdout3), resp.Stdout3)
	}
}
```
