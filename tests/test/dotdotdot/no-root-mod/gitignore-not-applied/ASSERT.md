## Expected
- Exit code 0.
- stderr contains mod_a, mod_b (non-gitignored modules found).
- stderr CONTAINS ign — gitignore is NOT applied because there is no git repository.
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
    for _, name := range []string{"mod_a", "mod_b", "ign"} {
        if !strings.Contains(resp.Stderr, name) {
            t.Fatalf("stderr missing %s (all dirs should be found when no git repo):\n%s", name, resp.Stderr)
        }
    }
    if !strings.Contains(resp.Stderr, "PASS") {
        t.Fatalf("stderr missing PASS:\n%s", resp.Stderr)
    }
}
```
