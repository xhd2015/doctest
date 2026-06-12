## Expected
- Exit code 0.
- stderr contains both pkg1 and pkg2 (sibling modules both found).
- stderr contains PASS.

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
    if !strings.Contains(resp.Stderr, "pkg1") {
        t.Fatalf("stderr missing pkg1:\n%s", resp.Stderr)
    }
    if !strings.Contains(resp.Stderr, "pkg2") {
        t.Fatalf("stderr missing pkg2:\n%s", resp.Stderr)
    }
    if !strings.Contains(resp.Stderr, "PASS") {
        t.Fatalf("stderr missing PASS:\n%s", resp.Stderr)
    }
}
```
