---
label: heavy
explanation: L3 product binary smoke for metrics --help argv/usage wiring
---

## Expected

- Exit code 0.
- Help text mentions each analyze subcommand: path, last, top, summary, show, prune.

## Exit Code

- 0

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("exit=%d stderr=%s stdout=%s", resp.ExitCode, resp.Stderr, resp.Stdout)
	}
	out := combinedOut(resp)
	// Help may be on stdout or stderr depending on CLI framework.
	for _, want := range []string{"path", "last", "top", "summary", "show", "prune"} {
		if !strings.Contains(out, want) {
			t.Fatalf("metrics help missing %q:\n%s", want, out)
		}
	}
}
```
