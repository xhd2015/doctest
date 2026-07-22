---
label: heavy
explanation: nested multi-arg thrice; third run uses -count=1
---

## Expected

- All three runs exit 0.
- Run2 (default warm) has **total Cached >= 2**.
- Run3 with `-count=1` has **total Cached == 0**.

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
	got2 := sumCachedCount(resp.Stdout2)
	if got2 < 2 {
		t.Fatalf("run2 warm multi-arg must total Cached >= 2 before testing -count; got %d\nstdout2:\n%s", got2, resp.Stdout2)
	}
	got3 := sumCachedCount(resp.Stdout3)
	if got3 != 0 {
		t.Fatalf("-count=1 on multi-arg must yield 0 Cached total; got %d\nstdout3:\n%s", got3, resp.Stdout3)
	}
}
```
