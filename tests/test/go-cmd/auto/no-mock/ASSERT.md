## Expected

- `Run` succeeds (`ExitCode` 0, empty error).
- `NeedsXgo` is false (no mock reachable).
- `ResolvedCmd` is exactly `go` (not `xgo`, not a path).

## Errors

- Classic TDD **RED** until `DetectXgoMockUsage` + `ResolveGoTestCmd` exist and
  auto maps needsXgo=false → `go`.

```go
import (
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run returned err: %v", err)
	}
	if resp == nil {
		t.Fatal("nil response")
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected exit 0 for auto+no-mock, got %d: %s", resp.ExitCode, errText(resp))
	}
	if resp.NeedsXgo {
		t.Fatalf("expected NeedsXgo=false for no-mock fixture, got true")
	}
	if resp.ResolvedCmd != "go" {
		t.Fatalf("auto+no-mock must resolve to go, got %q (err=%s)", resp.ResolvedCmd, errText(resp))
	}
}
```
