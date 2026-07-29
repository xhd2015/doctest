## Expected

- Non-zero `ExitCode`.
- Error text mentions `go-cmd`.
- Error indicates invalid value and/or lists allowed modes (`auto`, `xgo`, `go`)
  and/or includes the bad value `foo`.
- Must **not** succeed with `ResolvedCmd` set (parse-only path).

## Errors

- Classic TDD **RED** until `--go-cmd` is a known flag with value validation.
- Note: a bare `unrecognized flag: --go-cmd` without allowed-value guidance is
  **not** sufficient (still RED for product completeness once flag exists;
  before the flag exists, compile/parse may fail for other reasons).

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run returned err: %v", err)
	}
	if resp == nil {
		t.Fatal("nil response")
	}
	if resp.ExitCode == 0 {
		t.Fatalf("expected non-zero exit for --go-cmd=foo")
	}
	msg := errText(resp)
	lower := strings.ToLower(msg)
	if !strings.Contains(lower, "go-cmd") {
		t.Fatalf("error must mention go-cmd, got:\n%s", msg)
	}
	// Clear guidance: invalid value, allowed set, or the bad token.
	hasGuidance := strings.Contains(lower, "invalid") ||
		strings.Contains(lower, "unknown") ||
		strings.Contains(lower, "must be") ||
		strings.Contains(lower, "one of") ||
		(strings.Contains(lower, "auto") && strings.Contains(lower, "xgo") && strings.Contains(lower, "go")) ||
		strings.Contains(lower, "foo")
	if !hasGuidance {
		t.Fatalf("error must clarify invalid value or allowed auto|xgo|go, got:\n%s", msg)
	}
	// Reject silent accept.
	if resp.ResolvedCmd != "" {
		t.Fatalf("invalid go-cmd must not resolve a command, got %q", resp.ResolvedCmd)
	}
}
```
