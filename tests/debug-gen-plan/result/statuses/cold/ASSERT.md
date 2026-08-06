---
explanation: cold generate result tree + summary
---

## Expected

- Stderr includes `gen-plan: result` (same hierarchy shape as plan — at least
  `go.mod` / tree paths appear again under result).
- Summary line present: `summary:` with `new=`, `modified=`, `unchanged=`, `deleted=`.
- Cold gen: `new+modified` ≥ 1 (files written this run).
- Result tree lines use annotations `# new` / `# modified` / `# unchanged` / `# deleted`.
- `deleted` is 0 when prune removed nothing (typical cold empty gen).

## Errors

- None for GREEN after implement.

```go
import (
	"regexp"
	"strconv"
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run error: %v\nstderr:\n%s", err, resp.Stderr)
	}
	se := resp.Stderr
	if strings.Contains(se, `unknown key "gen-plan"`) {
		t.Fatalf("gen-plan not accepted:\n%s", se)
	}
	if !strings.Contains(se, "gen-plan: result") && !strings.Contains(se, "result") {
		t.Fatalf("expected gen-plan: result on stderr:\n%s", se)
	}
	// Tags on result lines (at least one write tag).
	if !strings.Contains(se, "# new") && !strings.Contains(se, "# modified") {
		t.Fatalf("expected # new or # modified annotation on result lines:\n%s", se)
	}
	// Summary: new=N modified=M unchanged=U deleted=K (flexible spacing)
	sumRe := regexp.MustCompile(`summary:.*new\s*=\s*(\d+).*modified\s*=\s*(\d+).*unchanged\s*=\s*(\d+).*deleted\s*=\s*(\d+)`)
	m := sumRe.FindStringSubmatch(se)
	if m == nil {
		// backward-compatible: modified= only (no new=)
		sumRe2 := regexp.MustCompile(`modified\s*=\s*(\d+).*unchanged\s*=\s*(\d+).*deleted\s*=\s*(\d+)`)
		m2 := sumRe2.FindStringSubmatch(se)
		if m2 == nil {
			t.Fatalf("expected summary new=/modified=/unchanged=/deleted= on stderr:\n%s", se)
		}
		if m2[1] == "0" {
			t.Fatalf("cold first run expected modified>=1, got %v\nstderr:\n%s", m2[0], se)
		}
		return
	}
	nw, _ := strconv.Atoi(m[1])
	mod, _ := strconv.Atoi(m[2])
	if nw+mod < 1 {
		t.Fatalf("cold first run expected new+modified>=1, got summary %v\nstderr:\n%s", m[0], se)
	}
}
```