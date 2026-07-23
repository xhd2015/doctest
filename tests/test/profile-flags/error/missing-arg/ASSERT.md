## Expected

- Non-zero exit code.
- Combined error mentions `cpuprofile` (flag name).
- Parse-time rejection of missing flag argument (not unknown-flag).

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	combined := resp.Stdout + resp.Stderr + resp.ParseErr
	lower := strings.ToLower(combined)

	if resp.ExitCode == 0 {
		t.Fatalf("expected non-zero exit for missing -cpuprofile value, got 0\ncombined:\n%s", combined)
	}
	// Classic TDD: "unknown flag" means the flag is not registered yet — that
	// is not the missing-argument contract. Require a parse error about the value.
	if strings.Contains(lower, "unknown flag") || strings.Contains(lower, "unknown option") ||
		strings.Contains(lower, "unrecognized flag") {
		t.Fatalf("got unknown-flag rejection (feature not registered); want missing-argument error for -cpuprofile:\n%s", combined)
	}
	if !strings.Contains(lower, "cpuprofile") {
		t.Fatalf("expected 'cpuprofile' in error output, got:\n%s", combined)
	}
	// Typical messages: "flag needs an argument", "missing argument", "requires an argument"
	if !strings.Contains(lower, "argument") && !strings.Contains(lower, "requires") &&
		!strings.Contains(lower, "expected") && !strings.Contains(lower, "need") {
		t.Fatalf("expected missing-argument style error for -cpuprofile, got:\n%s", combined)
	}
}
```
