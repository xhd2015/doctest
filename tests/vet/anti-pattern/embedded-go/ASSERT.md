## Expected
- The command fails with non-zero exit.
- stderr contains the embedded Go anti-pattern error.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("run failed unexpectedly: %v", err)
	}
	if resp.ExitCode == 0 {
		t.Fatalf("expected nonzero exit")
	}
	if !strings.Contains(resp.Stderr, "anti-pattern: raw Go code embedded in string literal") {
		t.Fatalf("stderr missing anti-pattern error:\n%s", resp.Stderr)
	}
}
```
