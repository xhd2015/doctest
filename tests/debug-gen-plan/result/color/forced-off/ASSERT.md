---
label: heavy
explanation: --no-color must strip ANSI from gen-plan result lines
---

## Expected

- Stderr includes gen-plan result markers / summary.
- Lines that are part of gen-plan output (containing `gen-plan:` or following
  result hierarchy / `summary:`) must not contain ESC (`\x1b`).
- Other stderr (unrelated product banners) may still be colored — assert
  focuses on lines with `gen-plan` or `summary: modified`.

## Errors

- RED until gen-plan result exists; when result lands without honoring
  --no-color, CSI on gen-plan lines fails this leaf.

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
	if !strings.Contains(se, "gen-plan:") {
		t.Fatalf("expected gen-plan markers on stderr:\n%s", se)
	}
	for _, line := range strings.Split(se, "\n") {
		// Focus on gen-plan-tagged lines and summary counts line.
		if strings.Contains(line, "gen-plan") || strings.Contains(line, "summary:") ||
			strings.Contains(line, "modified=") {
			if strings.Contains(line, "\x1b") {
				t.Fatalf("--no-color: gen-plan/summary line must not contain ANSI: %q", line)
			}
		}
	}
}
```
