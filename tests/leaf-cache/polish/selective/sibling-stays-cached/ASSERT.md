---
label: heavy
explanation: nested 3× doctest test on 2-leaf fixture with mid mutation
---

## Expected

- All three runs exit 0.
- Run2 warm: **Cached == 2** (both leaves skipped).
- Run3 after editing leaf_a: **Cached == 1** (leaf_b still hit; leaf_a re-ran).

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
	if resp.ExitCode3 != 0 {
		t.Fatalf("run3 exit %d\n%s\n%s", resp.ExitCode3, resp.Stdout3, resp.Stderr3)
	}
	c2 := cachedCount(resp.Stdout2)
	if c2 != 2 {
		t.Fatalf("run2 warm expected Cached==2, got %d\nstdout2:\n%s", c2, resp.Stdout2)
	}
	c3 := cachedCount(resp.Stdout3)
	if c3 != 1 {
		t.Fatalf("run3 after leaf_a edit expected Cached==1 (sibling only), got %d\nstdout3:\n%s", c3, resp.Stdout3)
	}
}
```
