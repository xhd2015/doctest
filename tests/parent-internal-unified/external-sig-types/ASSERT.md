## Expected

**Desired product behavior** (GREEN after fix; **RED** on current code):

- `RunErr` empty: expose for `internal/rules` compiles and subject suite runs.
- Combined stdout/stderr do **not** contain `undefined: model` (or other
  `undefined:` on external package names from expose facades).
- Prefer single suite package args when the go test line is present.

Today: `generateExposeSource` re-exports func wrappers with `model.Project` /
`model.FixResult` in the signature but does not import `example.com/app/model`
→ **RED**.

```go
import (
	"strings"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	t.Helper()
	_ = d
	_ = req
	if err != nil {
		t.Fatalf("unexpected harness error: %v", err)
	}
	if resp == nil {
		t.Fatal("nil response")
	}

	combined := resp.Stdout + "\n" + resp.Stderr + "\n" + resp.RunErr
	if strings.Contains(combined, "undefined: model") {
		t.Fatalf("expose must import external types used in signatures; got undefined: model\nstdout:\n%s\nstderr:\n%s\nRunErr=%s",
			resp.Stdout, resp.Stderr, resp.RunErr)
	}
	if resp.RunErr != "" {
		hint := ""
		if strings.Contains(combined, "undefined:") && strings.Contains(combined, "__doctest_internal_expose") {
			hint = " (expose missing imports for external types in signatures)"
		}
		t.Fatalf("want subject suite success with external-type internal API, got RunErr=%s%s\nstdout:\n%s\nstderr:\n%s",
			resp.RunErr, hint, resp.Stdout, resp.Stderr)
	}
	if len(resp.GoTestPackageArgs) > 1 {
		t.Fatalf("want single suite package, got pkgs=%v line=%q",
			resp.GoTestPackageArgs, resp.GoTestDisplayLine)
	}
}
```
