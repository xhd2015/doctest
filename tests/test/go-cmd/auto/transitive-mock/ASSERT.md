## Expected

- `Run` succeeds (`ExitCode` 0).
- `NeedsXgo` is true (mock reached only via helper, not direct entry import).
- `ResolvedCmd` is exactly `xgo`.

## Errors

- Classic TDD **RED** until transitive detection follows `runpkg` → `helper` →
  `github.com/xhd2015/xgo/runtime/mock` and auto resolves to `xgo`.
- Fail if detection only scans the entry package (false negative).

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
		t.Fatalf("expected exit 0 for auto+transitive-mock, got %d: %s", resp.ExitCode, errText(resp))
	}
	if !resp.NeedsXgo {
		t.Fatalf("expected NeedsXgo=true when helper imports %s (transitive from runpkg)", mockImportPath)
	}
	if resp.ResolvedCmd != "xgo" {
		t.Fatalf("auto+transitive-mock must resolve to xgo, got %q", resp.ResolvedCmd)
	}
}
```
