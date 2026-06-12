## Expected
- Non-zero exit code.
- stderr contains "read requirement file".

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if resp.ExitCode == 0 {
		t.Fatal("expected non-zero exit code for missing requirement file")
	}

	if !strings.Contains(resp.Stderr, "read requirement file") {
		t.Fatalf("stderr should contain 'read requirement file', got: %s", resp.Stderr)
	}
}
```
