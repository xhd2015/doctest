## Expected
- Nested fixture with leaf SETUP that has no `func Setup` is allowed (prose-only / no Setup).
- `doctest test` on the fixture exits 0.

```go
import (
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("run failed unexpectedly: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected zero exit for leaf SETUP without func Setup (now allowed), stdout:\n%s\nstderr:\n%s", resp.Stdout, resp.Stderr)
	}
}
```
