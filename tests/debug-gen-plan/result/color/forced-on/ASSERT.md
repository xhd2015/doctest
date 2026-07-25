---
label: heavy
explanation: force --color on gen-plan result statuses
---

## Expected

- Stderr includes `gen-plan: result`.
- Result-related stderr contains ANSI CSI sequences for status colors
  (green `\x1b[32m` / gray `\x1b[90m` or equivalent SGR used by product).
- Prefer matching via output assert `<ansi-color>` when lines are stable;
  here we check CSI presence near result section (structure may vary).

## Errors

- RED until colored result emit lands.

```go
import (
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
		t.Fatalf("expected gen-plan: result:\n%s", se)
	}
	// Must have some CSI (ESC[) in stderr for colored statuses.
	if !strings.Contains(se, "\x1b[") {
		t.Fatalf("expected ANSI CSI sequences with --color on gen-plan result, stderr:\n%s", se)
	}
	// Prefer green (modified) and/or gray (unchanged) SGR used by product color helpers.
	hasGreen := strings.Contains(se, "\x1b[32m") || strings.Contains(se, "\x1b[92m")
	hasGray := strings.Contains(se, "\x1b[90m") || strings.Contains(se, "\x1b[37m")
	if !hasGreen && !hasGray {
		t.Fatalf("expected green and/or gray SGR on result statuses, stderr:\n%s", se)
	}
}
```
