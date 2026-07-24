---
label: e2e, heavy
explanation: nested doctest test twice; compile fixture leaves
---

## Expected

- Both runs exit 0.
- Second run stdout summary has **Cached > 0** (programmatic leaf-cache skip).
- First run may be cold (`0 Cached`); not required to be zero if product pre-warms.

## Errors

- No harness error from Run.

## Exit Code

- Outer Run succeeds; nested exit codes both 0.

```go
import (
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v\nstderr1:\n%s\nstderr2:\n%s", err, resp.Stderr, resp.Stderr2)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("run1 exit %d\nstdout:\n%s\nstderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	if resp.ExitCode2 != 0 {
		t.Fatalf("run2 exit %d\nstdout:\n%s\nstderr:\n%s", resp.ExitCode2, resp.Stdout2, resp.Stderr2)
	}
	if !stdoutCachedPositive(resp.Stdout2) {
		t.Fatalf("warm second run expected Cached > 0; got count=%d\nstdout2:\n%s", cachedCount(resp.Stdout2), resp.Stdout2)
	}
}
```
