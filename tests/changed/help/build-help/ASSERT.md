---
explanation: L2 in-process CLI for build --help documents --changed
---

## Expected

- Exit code 0.
- stdout includes `Usage: doctest build` and `--changed`.

## Exit Code

0

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("exit code = %d\nstderr:\n%s", resp.ExitCode, resp.Stderr)
	}
	for _, want := range []string{"Usage: doctest build", "--changed"} {
		if !strings.Contains(resp.Stdout, want) {
			t.Fatalf("stdout missing %q:\n%s", want, resp.Stdout)
		}
	}
}
```
