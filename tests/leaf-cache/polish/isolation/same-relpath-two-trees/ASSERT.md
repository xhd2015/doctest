---
label: heavy
explanation: nested 3× doctest test across two twin trees
---

## Expected

- Run1 and run2 on treeA exit 0; run2 **Cached > 0**.
- Run3 on treeB exits 0 with **0 Cached** (no cross-tree pass hit).

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("run1 treeA exit %d\n%s\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	if resp.ExitCode2 != 0 {
		t.Fatalf("run2 treeA exit %d\n%s\n%s", resp.ExitCode2, resp.Stdout2, resp.Stderr2)
	}
	if !stdoutCachedPositive(resp.Stdout2) {
		t.Fatalf("run2 treeA warm expected Cached > 0; got %d\n%s", cachedCount(resp.Stdout2), resp.Stdout2)
	}
	if resp.ExitCode3 != 0 {
		t.Fatalf("run3 treeB exit %d\n%s\n%s", resp.ExitCode3, resp.Stdout3, resp.Stderr3)
	}
	if !stdoutCachedZero(resp.Stdout3) {
		t.Fatalf("treeB must not inherit treeA cache; got Cached=%d\nstdout3:\n%s", cachedCount(resp.Stdout3), resp.Stdout3)
	}
}
```
