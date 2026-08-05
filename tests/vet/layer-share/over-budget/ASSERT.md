## Expected

- Non-zero exit.
- Error text identifies the tree and L3 over-budget (stable substrings for implementer):
  - root path (Args dir) or clear tree identifier
  - `L3` and share/pct or fraction
  - `10%` or `max 10`
- Example shape (wording may vary):  
  `<path>: L3 share 20.0% (2/10 leaves labeled e2e) exceeds max 10%`

## Exit Code

- non-zero

```go
import (
	"path/filepath"
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	if err != nil {
		t.Fatalf("run failed unexpectedly: %v", err)
	}
	if resp.ExitCode == 0 {
		t.Fatalf("expected nonzero exit when L3 share over budget (10 leaves, 2 e2e)")
	}
	msg := resp.Stderr
	if strings.TrimSpace(msg) == "" {
		msg = resp.Stdout
	}
	if !strings.Contains(msg, "L3") {
		t.Fatalf("error must mention L3:\n%s", msg)
	}
	if !strings.Contains(msg, "10%") && !strings.Contains(msg, "max 10") {
		t.Fatalf("error must mention budget bound (10%% or max 10):\n%s", msg)
	}
	// Share signal: percent figure and/or fraction and/or word "share".
	hasShareSignal := strings.Contains(strings.ToLower(msg), "share") ||
		strings.Contains(msg, "%") ||
		strings.Contains(msg, "/")
	if !hasShareSignal {
		t.Fatalf("error must include share/pct or fraction:\n%s", msg)
	}
	dir := ""
	if len(req.Args) > 0 {
		dir = req.Args[len(req.Args)-1]
	}
	if dir != "" {
		base := filepath.Base(dir)
		if !strings.Contains(msg, dir) && !strings.Contains(msg, base) {
			t.Fatalf("error must identify tree root %q (or basename %q):\n%s", dir, base, msg)
		}
	}
}
```
