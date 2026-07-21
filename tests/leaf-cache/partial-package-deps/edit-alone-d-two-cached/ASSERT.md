---
label: heavy
explanation: nested 2× doctest; alone/d edit → 2 Cached
---

## Expected

- Run1 (`-count=1`) exit 0 and **0 Cached**.
- After editing `alone/d`, run2 exit 0 and **exactly 2 Cached**
  (leaf-ab-1 + leaf-ab-2 warm; leaf-d re-runs).

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("run1 exit %d\nstdout:\n%s\nstderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	if !stdoutCachedZero(resp.Stdout) {
		t.Fatalf("run1 (-count=1) expected 0 Cached; got %d\n%s", cachedCount(resp.Stdout), resp.Stdout)
	}
	if resp.ExitCode2 != 0 {
		t.Fatalf("run2 exit %d\nstdout:\n%s\nstderr:\n%s", resp.ExitCode2, resp.Stdout2, resp.Stderr2)
	}
	got := cachedCount(resp.Stdout2)
	if got != 2 {
		t.Fatalf("after alone/d edit expected 2 Cached (shared leaves warm); got %d\nstdout2:\n%s", got, resp.Stdout2)
	}
}
```
