---
label: heavy
---

## Expected
- `build.Test` succeeds (one passing leaf).
- Captured stdout contains the summary line but no ANSI escape sequences.

## Exit Code
- `err` is nil.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if resp.TestErr != nil {
		t.Fatalf("expected build.Test to succeed, got: %v", resp.TestErr)
	}
	if !strings.Contains(resp.Summary, "1 Pass") {
		t.Fatalf("expected summary with 1 Pass, got:\n%s", resp.Output)
	}
	if containsANSI(resp.Output) {
		t.Fatalf("ColorAuto on opts.Stdout buffer (non-TTY) must not emit ANSI, got:\n%s", resp.Output)
	}
}
```