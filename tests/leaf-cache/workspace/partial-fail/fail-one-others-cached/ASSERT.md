---
label: heavy
explanation: nested multi-tree workspace; partial fail then warm pass leaf
---

## Expected

- Both nested runs exit **non-zero** (fail leaf still present).
- Run2 summary has **Cached >= 1** (the previously passing leaf is warm-skipped).
- Run1 may have 0 Cached (cold).

## Errors

- Outer Run itself should not error (captures exit codes in Response).

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected harness error: %v", err)
	}
	if resp.ExitCode == 0 {
		t.Fatalf("run1 expected failure exit, got 0\nstdout:\n%s\nstderr:\n%s", resp.Stdout, resp.Stderr)
	}
	if resp.ExitCode2 == 0 {
		t.Fatalf("run2 expected failure exit, got 0\nstdout:\n%s\nstderr:\n%s", resp.Stdout2, resp.Stderr2)
	}
	got := cachedCount(resp.Stdout2)
	if got < 1 {
		t.Fatalf("run2 expected Cached >= 1 for the warm pass leaf; got %d\nstdout2:\n%s", got, resp.Stdout2)
	}
}
```
