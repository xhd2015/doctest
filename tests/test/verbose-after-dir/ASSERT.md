## Expected
- The command succeeds.
- The verbose flag is forwarded to the underlying test-case-tree runner.

```go
import (
    "strings"
    "testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if err != nil {
        t.Fatalf("run failed unexpectedly: %v", err)
    }
    if resp.ExitCode != 0 {
        t.Fatalf("expected zero exit, got %d\nstdout:\n%s\nstderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
    }
    if strings.Contains(resp.Stderr, "test requires <dir>") {
        t.Fatalf("verbose flag after dir was treated as an extra positional arg:\n%s", resp.Stderr)
    }
    if !strings.Contains(resp.Stderr, "TestGeneratedCaseHappyPath") {
        t.Fatalf("expected verbose test output to include generated test name, stderr:\n%s", resp.Stderr)
    }
}
```
