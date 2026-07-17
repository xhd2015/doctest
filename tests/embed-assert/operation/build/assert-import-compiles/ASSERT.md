---
label: heavy
---

## Expected

- `doctest build` exits 0.
- Generated test package compiles with assert module resolved via cache replace.

## Exit Code

- Exit code 0.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("doctest subprocess failed: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected doctest build with assert to succeed, got exit %d\nstdout:\n%s\nstderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
}
```