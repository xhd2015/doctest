---
label: heavy
explanation: nested doctest with --color on pass+fail fixture
---

## Expected

- Nested run exits non-zero (has a fail leaf).
- Progress region contains **>= 1** red ANSI dot (`\x1b[31m.\x1b[0m`).

## Errors

- No harness error from Run.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected harness error: %v", err)
	}
	if resp.ExitCode == 0 {
		t.Fatalf("expected fail exit, got 0\nstdout:\n%s", resp.Stdout)
	}
	red := countRedProgressDots(resp.Stdout)
	if red < 1 {
		t.Fatalf("fail progress dot must be red under --color; red=%d\nprogress:\n%q\nstdout:\n%s",
			red, progressPrefix(resp.Stdout), resp.Stdout)
	}
}
```
