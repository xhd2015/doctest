---
label: heavy
explanation: nested doctest test thrice; third run uses -a
---

## Expected

- All three runs exit 0.
- Run2 warm: **Cached > 0**.
- Run3 with `-a`: **0 Cached**.

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
	if !stdoutCachedPositive(resp.Stdout2) {
		t.Fatalf("run2 warm must Cached > 0 before testing -a; got %d\nstdout2:\n%s", cachedCount(resp.Stdout2), resp.Stdout2)
	}
	if !stdoutCachedZero(resp.Stdout3) {
		t.Fatalf("-a must yield 0 Cached; got %d\nstdout3:\n%s", cachedCount(resp.Stdout3), resp.Stdout3)
	}
}
```
