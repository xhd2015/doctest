---
explanation: nested doctest twice with --color on warm run; multi-leaf fixture
---

## Expected

- Both runs exit 0.
- Run2 summary **Cached >= 2**.
- Progress region on run2 contains **>= 2** grey ANSI dots
  (`\x1b[90m.\x1b[0m`) — one per warm leaf-cache skip.

## Expected Output

Progress region (before the suite `(N Run, …)` line) must include grey dots for
cached skips. Exact template is timing-sensitive; assert uses count helpers.

## Errors

- No harness error from Run.

## Exit Code

- Nested exits both 0.

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v\nstderr1:\n%s\nstderr2:\n%s", err, resp.Stderr, resp.Stderr2)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("run1 exit %d\nstdout:\n%s\nstderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	if resp.ExitCode2 != 0 {
		t.Fatalf("run2 exit %d\nstdout:\n%s\nstderr:\n%s", resp.ExitCode2, resp.Stdout2, resp.Stderr2)
	}
	if cachedCount(resp.Stdout2) < 2 {
		t.Fatalf("warm --color run expected Cached >= 2; got %d\nstdout2:\n%s",
			cachedCount(resp.Stdout2), resp.Stdout2)
	}
	grey := countGrayProgressDots(resp.Stdout2)
	if grey < 2 {
		t.Fatalf("warm leaf-cache skips must print grey progress dots; got grey=%d Cached=%d\n"+
			"progress:\n%q\nfull stdout2:\n%s",
			grey, cachedCount(resp.Stdout2), progressPrefix(resp.Stdout2), resp.Stdout2)
	}
}
```
