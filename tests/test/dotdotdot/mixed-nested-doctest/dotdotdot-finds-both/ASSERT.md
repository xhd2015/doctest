## Expected
- Exit code 0.
- `./ancestor/leaf/...` finds both the leaf itself (via ancestor tree) and the nested sub2 DOCTEST root.
- stderr contains "ancestor" (the ancestor doctest root is resolved).
- stderr contains "sub2" (the nested DOCTEST root is discovered).
- stderr contains PASS.

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
    if !strings.Contains(resp.Stderr, "ancestor") {
        t.Fatalf("stderr missing ancestor:\n%s", resp.Stderr)
    }
    if !strings.Contains(resp.Stderr, "sub2") {
        t.Fatalf("stderr missing sub2:\n%s", resp.Stderr)
    }
    if !strings.Contains(resp.Stdout, "PASS") {
        t.Fatalf("stdout missing PASS:\n%s", resp.Stdout)
    }
}
```
