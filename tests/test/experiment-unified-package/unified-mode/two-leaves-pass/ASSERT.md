---
label: heavy
---

## Expected

- Suite run succeeds (`RunErr` empty) under default unified generation.
- Fixture leaves `a` and `b` both executed successfully (implied by suite success).

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected harness error: %v", err)
	}
	if resp.RunErr != "" {
		t.Fatalf("unified-mode 2-leaf RunTest failed: %s\nstdout:\n%s\nstderr:\n%s\ngen=%s",
			resp.RunErr, resp.Stdout, resp.Stderr, resp.GenDir)
	}
}
```
