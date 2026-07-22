---
label: heavy
explanation: nested 3× doctest; warm grey then -count=1 --color
---

## Expected

- Run1 and run2 exit 0; run2 Cached >= 2 and grey progress dots >= 2
  (warm path armed).
- Run3 with `-count=1 --color` exits 0, **0 Cached**, and **0** grey
  progress dots (no leaf-cache skip dots).

## Errors

- No harness error from Run.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v\nstderr:\n%s\n%s\n%s", err, resp.Stderr, resp.Stderr2, resp.Stderr3)
	}
	if resp.ExitCode != 0 || resp.ExitCode2 != 0 {
		t.Fatalf("run1/run2 must pass; exit1=%d exit2=%d\nstdout2:\n%s",
			resp.ExitCode, resp.ExitCode2, resp.Stdout2)
	}
	if cachedCount(resp.Stdout2) < 2 || countGrayProgressDots(resp.Stdout2) < 2 {
		t.Fatalf("precondition: warm --color must Cached>=2 and grey>=2; Cached=%d grey=%d\nstdout2:\n%s",
			cachedCount(resp.Stdout2), countGrayProgressDots(resp.Stdout2), resp.Stdout2)
	}
	if resp.ExitCode3 != 0 {
		t.Fatalf("run3 -count=1 exit %d\nstdout3:\n%s\nstderr3:\n%s",
			resp.ExitCode3, resp.Stdout3, resp.Stderr3)
	}
	if !stdoutCachedZero(resp.Stdout3) {
		t.Fatalf("-count=1 must yield 0 Cached; got %d\nstdout3:\n%s",
			cachedCount(resp.Stdout3), resp.Stdout3)
	}
	grey := countGrayProgressDots(resp.Stdout3)
	if grey != 0 {
		t.Fatalf("-count=1 must not emit grey leaf-cache skip dots; grey=%d\nprogress:\n%q\nstdout3:\n%s",
			grey, progressPrefix(resp.Stdout3), resp.Stdout3)
	}
}
```
