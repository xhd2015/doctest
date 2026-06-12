## Expected
- Exit code 0 (no tests found is not an error).
- stderr contains "no tests found".

## Exit Code
- Exit code 0.

```go
import (
    "strings"
    "testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if err != nil {
        t.Fatalf("run failed: %v", err)
    }
    if resp.ExitCode != 0 {
        t.Fatalf("exit code = %d, stderr:\n%s", resp.ExitCode, resp.Stderr)
    }
    if !strings.Contains(resp.Stderr, "no tests found") {
        t.Fatalf("stderr missing 'no tests found':\n%s", resp.Stderr)
    }
}
```
