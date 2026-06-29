## Expected

- Non-zero exit, parse error on stderr, no PASS/FAIL summary.

## Exit Code

non-zero

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}
	if resp.ExitCode == 0 {
		t.Fatalf("expected nonzero exit\nstderr:\n%s", resp.Stderr)
	}
	if !strings.Contains(resp.Stderr, "parse") && !strings.Contains(resp.Stderr, "syntax") {
		t.Fatalf("stderr must report expression parse/syntax error, not generic flag error:\n%s", resp.Stderr)
	}
	assertNoResultSummary(t, resp.Stdout)
}
```