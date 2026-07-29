## Expected

- `ExitCode` 0.
- `NeedsXgo` is true (detection still reports mock usage).
- `ResolvedCmd` is exactly `go` (force overrides auto choice of xgo).

## Errors

- Classic TDD **RED** until force-go path of `ResolveGoTestCmd` lands.
- Documented product note: mocks may fail at runtime under plain `go test`;
  CLI still honors `--go-cmd=go`.

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run returned err: %v", err)
	}
	if resp == nil {
		t.Fatal("nil response")
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected exit 0 for force-go, got %d: %s", resp.ExitCode, errText(resp))
	}
	if !resp.NeedsXgo {
		t.Fatalf("detection should still report NeedsXgo=true for transitive mock fixture")
	}
	if resp.ResolvedCmd != "go" {
		t.Fatalf("force --go-cmd=go must resolve to go despite mock, got %q", resp.ResolvedCmd)
	}
}
```
