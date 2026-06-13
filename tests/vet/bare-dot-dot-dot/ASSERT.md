## Expected
- The command fails with non-zero exit code.
- stderr contains `bare '...' pattern is not supported`.

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
		t.Fatalf("expected nonzero exit, stdout:\n%s", resp.Stdout)
	}
	if !strings.Contains(resp.Stderr, "bare '...' pattern is not supported") {
		t.Fatalf("stderr missing expected error:\n%s", resp.Stderr)
	}
}
```
