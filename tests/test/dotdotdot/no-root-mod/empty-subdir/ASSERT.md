---
label: heavy
---

## Expected
- Exit code 0 (no tests is not an error).
- stderr contains "no tests".

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
    if !strings.Contains(resp.Stderr, "no tests") {
        t.Fatalf("stderr missing 'no tests':\n%s", resp.Stderr)
    }
}
```
