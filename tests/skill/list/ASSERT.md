## Expected
- The command succeeds.
- stdout lists both available spec skills.

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
    for _, want := range []string{"doc-spec", "code-spec", "tdd", "tdd-lite", "reproduce", "review", "implementer"} {
        if !strings.Contains(resp.Stdout, want) {
            t.Fatalf("stdout missing %q:\n%s", want, resp.Stdout)
        }
    }
}
```

