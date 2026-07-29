## Expected

- Non-zero `ExitCode`.
- `ResolvedCmd` is `xgo` (policy chose xgo before lookup failed), **or** error
  path may leave it set — either is fine if error text is clear.
- Error text mentions `xgo` and that it was not found / missing from PATH
  (actionable; not a generic exec failure).

## Errors

- Classic TDD **RED** until `EnsureGoTestCmdAvailable` fails clearly for missing xgo.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run returned err: %v (expect error via resp.ExitCode/ErrMsg)", err)
	}
	if resp == nil {
		t.Fatal("nil response")
	}
	if resp.ExitCode == 0 {
		t.Fatalf("expected non-zero exit when xgo missing from search PATH")
	}
	msg := strings.ToLower(errText(resp))
	if !strings.Contains(msg, "xgo") {
		t.Fatalf("error must mention xgo, got:\n%s", errText(resp))
	}
	// Actionable: not found / PATH
	if !strings.Contains(msg, "not found") && !strings.Contains(msg, "path") {
		t.Fatalf("error must be actionable (not found / PATH), got:\n%s", errText(resp))
	}
	// Prefer resolved name visible when lookup fails after resolve.
	if resp.ResolvedCmd != "" && resp.ResolvedCmd != "xgo" {
		t.Fatalf("if ResolvedCmd is set, expected xgo, got %q", resp.ResolvedCmd)
	}
}
```
