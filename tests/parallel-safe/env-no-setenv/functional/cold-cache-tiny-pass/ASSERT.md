---
label: heavy
explanation: nested doctest test --cold-cache generate + go test smoke
---

## Expected

- Subprocess exits 0.
- Stderr (or stdout) announces cold-cache mode and mentions GOCACHE isolation.
- GREEN before and after P1: after fix, isolation must still work via
  `opts.GoCache` → child `cmd.Env` (no process `os.Setenv("GOCACHE")`).

## Exit Code

- 0

```go
import (
	"strings"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	if err != nil {
		t.Fatalf("run failed: %v\nstderr:\n%s", err, resp.Stderr)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected exit 0 for cold-cache tiny fixture, got %d\nstdout:\n%s\nstderr:\n%s",
			resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	combined := strings.ToLower(resp.Stderr + "\n" + resp.Stdout)
	if !strings.Contains(combined, "cold-cache") && !strings.Contains(combined, "cold cache") {
		t.Fatalf("expected cold-cache announcement on stderr:\n%s", resp.Stderr)
	}
	if !strings.Contains(combined, "gocache") && !strings.Contains(combined, "go cache") {
		t.Fatalf("expected cold-cache announcement to mention GOCACHE:\n%s", resp.Stderr)
	}
}
```
