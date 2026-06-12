## Expected
- Exit code 0.
- stderr contains mod_a, mod_b, src (non-gitignored modules found).
- stderr does NOT contain ign_a, ign_b, foo_test, bar_test (gitignored modules skipped).
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
    for _, name := range []string{"mod_a", "mod_b", "src"} {
        if !strings.Contains(resp.Stderr, name) {
            t.Fatalf("stderr missing %s:\n%s", name, resp.Stderr)
        }
    }
    for _, name := range []string{"ign_a", "ign_b", "foo_test", "bar_test"} {
        if strings.Contains(resp.Stderr, name) {
            t.Fatalf("stderr should not contain gitignored %s:\n%s", name, resp.Stderr)
        }
    }
    if !strings.Contains(resp.Stderr, "PASS") {
        t.Fatalf("stderr missing PASS:\n%s", resp.Stderr)
    }
}
```
