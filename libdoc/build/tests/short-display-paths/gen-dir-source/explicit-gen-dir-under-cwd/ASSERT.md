---
label: heavy
---

## Expected

- `build.Test` succeeds.
- `resp.HeaderLine` is `doctest: tests/feature (1 tests)`.
- `resp.CdLine` uses `_gen/` for the explicit gen dir (not an absolute temp path).

## Exit Code

- `err` is nil.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if resp.TestErr != nil {
		t.Fatalf("expected build.Test to succeed, got: %v\nstderr:\n%s", resp.TestErr, resp.Stderr)
	}
	if resp.HeaderLine != "doctest: tests/feature (1 tests)" {
		t.Fatalf("expected compact header, got %q\nstderr:\n%s", resp.HeaderLine, resp.Stderr)
	}
	if resp.CdLine == "" {
		t.Fatalf("missing cd line in stderr:\n%s", resp.Stderr)
	}
	if !strings.Contains(resp.CdLine, "_gen/") {
		t.Fatalf("cd line must use _gen/, got: %s", resp.CdLine)
	}
	if strings.Contains(resp.CdLine, "/var/") || strings.Contains(resp.CdLine, "\\Temp\\") {
		t.Fatalf("explicit gen dir must not show temp absolute path: %s", resp.CdLine)
	}
}
```