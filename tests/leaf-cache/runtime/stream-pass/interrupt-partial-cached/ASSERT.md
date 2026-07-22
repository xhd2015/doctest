---
label: heavy
explanation: nested doctest interrupt mid-run then warm re-run; compile multi-leaf fixture
---

## Expected

- Run1 is interrupted after progress dots (non-zero exit is OK; harness must
  have delivered SIGINT).
- After unhang, run2 exits 0.
- Run2 summary has **Cached >= 1** — leaves that passed before SIGINT were
  stream-PutPass'd and warm-skipped.

## Errors

- No harness error from Run (SIGINT path returns captured exit, not Run error).

## Exit Code

- Outer Run succeeds; run2 nested exit 0; run1 may be non-zero.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected harness error: %v\nstderr1:\n%s\nstderr2:\n%s", err, resp.Stderr, resp.Stderr2)
	}
	if resp.ExitCode2 != 0 {
		t.Fatalf("run2 after unhang expected exit 0, got %d\nstdout2:\n%s\nstderr2:\n%s",
			resp.ExitCode2, resp.Stdout2, resp.Stderr2)
	}
	got := cachedCount(resp.Stdout2)
	if got < 1 {
		t.Fatalf("stream PutPass: after mid-run interrupt, run2 expected Cached >= 1; got %d\n"+
			"stdout1 (interrupted):\n%s\nstdout2:\n%s",
			got, resp.Stdout, resp.Stdout2)
	}
}
```
