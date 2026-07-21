---
label: heavy
explanation: nested 2× doctest; shared/a edit → 1 Cached
---

## Expected

- Run1 (`-count=1`) exit 0 and **0 Cached**.
- After editing `shared/a`, run2 exit 0 and **exactly 1 Cached**
  (leaf-d warm; leaf-ab-1 + leaf-ab-2 re-run).

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
	if got != 1 {
		t.Fatalf("after shared/a edit expected 1 Cached (leaf-d only); got %d\nstdout2:\n%s", got, resp.Stdout2)
	}
}
```
