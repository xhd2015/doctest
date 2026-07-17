---
label: heavy
---

## Expected
- The command succeeds.
- stderr shows the go test command with `-timeout=45s`.

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
        t.Fatalf("expected exit code 0, got %d\nstdout:\n%s\nstderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
    }
    if !strings.Contains(resp.Stderr, "-timeout=45s") {
        t.Fatalf("expected stderr to contain '-timeout=45s', got:\n%s", resp.Stderr)
    }
}
```