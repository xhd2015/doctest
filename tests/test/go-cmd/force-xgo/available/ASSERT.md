## Expected

- `ExitCode` 0.
- `ResolvedCmd` is exactly `xgo` despite `NeedsXgo=false`.

## Errors

- Classic TDD **RED** until force-xgo path of `ResolveGoTestCmd` lands.

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
		t.Fatalf("expected exit 0 for force-xgo available, got %d: %s", resp.ExitCode, errText(resp))
	}
	if resp.ResolvedCmd != "xgo" {
		t.Fatalf("force --go-cmd=xgo must resolve to xgo even without mock, got %q", resp.ResolvedCmd)
	}
}
```
