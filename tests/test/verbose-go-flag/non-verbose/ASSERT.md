---
label: heavy
---

## Expected
- The command succeeds.
- stderr shows the go test command WITHOUT `-v`.

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
    if !strings.Contains(resp.Stderr, "go test ") {
        t.Fatalf("expected stderr to contain go test command preview, got:\n%s", resp.Stderr)
    }
    if strings.Contains(resp.Stderr, " -v") {
        t.Fatalf("expected stderr NOT to contain '-v', got:\n%s", resp.Stderr)
    }
}
```
