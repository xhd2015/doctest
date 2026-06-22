## Expected
- `build.Test` fails (one failing package).
- First dot (pass) is plain; second dot (fail) has ANSI red.
- Summary has green `1 Pass`, red `1 Fail`, plain `2 Run`.

## Errors
- `resp.TestErr` is non-nil.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run returned unexpected error: %v", err)
	}
	if resp.TestErr == nil {
		t.Fatal("expected build.Test to fail for mixed pass/fail")
	}

	plainDots := stripANSI(resp.Dots)
	if plainDots != ".." {
		t.Fatalf("expected two dots, got plain %q raw %q", plainDots, resp.Dots)
	}
	if strings.Count(resp.Dots, "\x1b[31m") != 1 {
		t.Fatalf("expected exactly one red fail dot, dots:\n%q", resp.Dots)
	}

	if !metricIsColored(resp.Summary, "1 Pass") {
		t.Fatalf("expected green 1 Pass, summary:\n%s", resp.Summary)
	}
	if !failFieldIsColored(resp.Summary) {
		t.Fatalf("expected red 1 Fail field, summary:\n%s", resp.Summary)
	}
	if !metricIsPlain(resp.Summary, "2 Run") {
		t.Fatalf("expected plain 2 Run, summary:\n%s", resp.Summary)
	}
}
```