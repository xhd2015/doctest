---
label: heavy
explanation: warm single tree-a then workspace; guard bare-path false skip
---

## Expected

- Run1 (single-tree tree-a) exits 0.
- Run2 (workspace `/...`) exits **non-zero** — tree-b fail body must execute.
- Run2 **Cached == 1** (only tree-a is a warm GetPass hit).

## Errors

- Exit 0 on run2 with Cached >= 2 indicates cross-tree false skip (bare relpath).

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("run1 tree-a warm exit %d\nstdout:\n%s\nstderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	if resp.ExitCode2 == 0 {
		t.Fatalf("run2 workspace must fail (tree-b body); exit 0 suggests false warm skip\nstdout2:\n%s\nstderr2:\n%s", resp.Stdout2, resp.Stderr2)
	}
	got := cachedCount(resp.Stdout2)
	if got != 1 {
		t.Fatalf("run2 expected Cached == 1 (tree-a only); got %d\nstdout2:\n%s", got, resp.Stdout2)
	}
}
```
