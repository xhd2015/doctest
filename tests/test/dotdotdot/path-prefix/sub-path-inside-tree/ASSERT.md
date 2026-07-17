---
label: heavy
---

## Expected
- Exit code 0.
- `./alpha/simple/...` resolves the doctest root by walking up and runs tests scoped to `alpha/simple/`.
- stderr contains "alpha" (the doctest root is found).
- stderr does NOT contain "beta" (sibling doctest tree excluded).
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
    if !strings.Contains(resp.Stderr, "alpha") {
        t.Fatalf("stderr missing alpha:\n%s", resp.Stderr)
    }
    if strings.Contains(resp.Stderr, "beta") {
        t.Fatalf("stderr should not contain beta:\n%s", resp.Stderr)
    }
    if !strings.Contains(resp.Stdout, "PASS") {
        t.Fatalf("stdout missing PASS:\n%s", resp.Stdout)
    }
}
```
