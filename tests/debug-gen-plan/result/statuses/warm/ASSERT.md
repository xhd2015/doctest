---
explanation: two generates same gen-dir; warm unchanged summary
---

## Expected

- Both runs emit `gen-plan: result` with a summary line.
- Second run: `unchanged` ≥ first-run `unchanged`, and typically second
  `unchanged` ≥ 1 (hash-hit skips rewrite).
- Or equivalently: second `modified` ≤ first `modified` with second
  `unchanged` > 0.

## Errors

- RED until result summary and write tracking land.

```go
import (
	"regexp"
	"strconv"
	"strings"
	"testing"
)

func parseSummary(stderr string) (mod, unc, del int, ok bool) {
	re := regexp.MustCompile(`modified\s*=\s*(\d+).*unchanged\s*=\s*(\d+).*deleted\s*=\s*(\d+)`)
	m := re.FindStringSubmatch(stderr)
	if m == nil {
		return 0, 0, 0, false
	}
	mod, _ = strconv.Atoi(m[1])
	unc, _ = strconv.Atoi(m[2])
	del, _ = strconv.Atoi(m[3])
	return mod, unc, del, true
}

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run error: %v\nstderr1:\n%s\nstderr2:\n%s", err, resp.Stderr, resp.SecondStderr)
	}
	for i, se := range []string{resp.Stderr, resp.SecondStderr} {
		if strings.Contains(se, `unknown key "gen-plan"`) {
			t.Fatalf("run %d: gen-plan not accepted:\n%s", i+1, se)
		}
	}
	mod1, unc1, _, ok1 := parseSummary(resp.Stderr)
	mod2, unc2, _, ok2 := parseSummary(resp.SecondStderr)
	if !ok1 {
		t.Fatalf("first run missing summary:\n%s", resp.Stderr)
	}
	if !ok2 {
		t.Fatalf("second run missing summary:\n%s", resp.SecondStderr)
	}
	// Warm: more unchanged and/or fewer modified than cold.
	if unc2 < unc1 {
		t.Fatalf("warm unchanged should not drop: first modified=%d unchanged=%d; second modified=%d unchanged=%d", mod1, unc1, mod2, unc2)
	}
	if unc2 == 0 && mod2 >= mod1 && mod1 > 0 {
		t.Fatalf("warm second run expected some unchanged (or fewer modified); first m=%d u=%d; second m=%d u=%d\nstderr2:\n%s",
			mod1, unc1, mod2, unc2, resp.SecondStderr)
	}
	if unc2 == 0 {
		// Allow if second fully rewrites but then require mod2 < mod1 as weak signal;
		// preferred is unc2 > 0.
		t.Fatalf("warm second run expected unchanged>=1, got modified=%d unchanged=%d\n%s", mod2, unc2, resp.SecondStderr)
	}
	if !strings.Contains(resp.SecondStderr, "# unchanged") {
		t.Fatalf("warm second run expected # unchanged annotations:\n%s", resp.SecondStderr)
	}
}
```
