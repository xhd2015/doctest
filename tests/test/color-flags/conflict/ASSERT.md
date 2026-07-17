---
label: heavy
---

## Expected
- The command rejects `--color` and `--no-color` together with an explicit conflict error.
- Exit code is non-zero.
- Error must not be a generic "unrecognized flag" message.

## Errors
- `resp.ExitCode` is non-zero.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	combined := resp.Stdout + resp.Stderr
	lower := strings.ToLower(combined)

	if strings.Contains(lower, "unrecognized flag") {
		t.Fatalf("expected mutual-exclusion error, not unrecognized flag:\n%s", combined)
	}

	conflict := strings.Contains(lower, "conflict") ||
		strings.Contains(lower, "mutually exclusive") ||
		(strings.Contains(lower, "--color") && strings.Contains(lower, "--no-color"))
	if !conflict {
		t.Fatalf("expected --color/--no-color conflict error, got:\n%s", combined)
	}

	if resp.ExitCode == 0 {
		t.Fatalf("expected non-zero exit for conflicting color flags, got 0\ncombined:\n%s", combined)
	}
}
```