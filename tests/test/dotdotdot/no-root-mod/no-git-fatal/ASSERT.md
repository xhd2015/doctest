## Expected
- Exit code 0.
- stderr does NOT contain "fatal: not a git repository".
- stderr contains mod_a, mod_b (all modules found).
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
    // Verify fatal git message does NOT appear
    if strings.Contains(resp.Stderr, "fatal: not a git repository") {
        t.Fatalf("stderr contains 'fatal: not a git repository' but should not:\n%s", resp.Stderr)
    }
    for _, name := range []string{"mod_a", "mod_b"} {
        if !strings.Contains(resp.Stderr, name) {
            t.Fatalf("stderr missing %s:\n%s", name, resp.Stderr)
        }
    }
    if !strings.Contains(resp.Stdout, "PASS") {
        t.Fatalf("stdout missing PASS:\n%s", resp.Stdout)
    }
}
```
