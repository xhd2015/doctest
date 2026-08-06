---
label: e2e
explanation: nested doctest test twice on failing fixture
---

## Expected

- Both nested runs exit **non-zero**.
- Both report **0 Cached** (fail never stored; no skip).

## Errors

- Outer Run itself should not error (captures exit codes in Response).

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected harness error: %v", err)
	}
	if resp.ExitCode == 0 {
		t.Fatalf("run1 expected failure exit, got 0\nstdout:\n%s", resp.Stdout)
	}
	if resp.ExitCode2 == 0 {
		t.Fatalf("run2 expected failure exit, got 0\nstdout:\n%s", resp.Stdout2)
	}
	if !stdoutCachedZero(resp.Stdout) {
		t.Fatalf("run1 must not Cached a fail; count=%d\n%s", cachedCount(resp.Stdout), resp.Stdout)
	}
	if !stdoutCachedZero(resp.Stdout2) {
		t.Fatalf("run2 must not Cached a fail; count=%d\n%s", cachedCount(resp.Stdout2), resp.Stdout2)
	}
}
```
