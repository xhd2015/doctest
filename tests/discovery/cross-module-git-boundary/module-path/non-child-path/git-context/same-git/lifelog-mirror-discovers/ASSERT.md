## Expected
- Exit code 0.
- `skill-cli` tests are discovered.
- No `warning:` on stderr.
- No `no tests`.

## Side Effects
- None beyond running nested doctest tree.

## Exit Code
- 0

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
    if strings.Contains(resp.Stderr, "warning:") {
        t.Fatalf("stderr should not contain warning (lifelog mirror, same git should discover):\n%s", resp.Stderr)
    }
    if !strings.Contains(resp.Stderr, "skill-cli") {
        t.Fatalf("stderr missing skill-cli (lifelog mirror should discover lifelog-cli tests):\n%s", resp.Stderr)
    }
    if strings.Contains(resp.Stderr, "no tests") {
        t.Fatalf("stderr should not contain 'no tests':\n%s", resp.Stderr)
    }
    if !strings.Contains(resp.Stdout, "PASS") {
        t.Fatalf("stdout missing PASS:\n%s", resp.Stdout)
    }
}
```
