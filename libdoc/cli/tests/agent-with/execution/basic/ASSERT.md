## Expected
- Exit code 0.
- Stdout contains `hello`.

## Side Effects
- The child process `echo hello` runs with `DOCTEST_SUBAGENT_AGENT_RUNNER=opencode` in its environment.

## Exit Code
- 0.

```go
import (
    "strings"
    "testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if resp.Err != nil {
        t.Fatalf("unexpected error: %v", resp.Err)
    }
    if resp.ExitCode != 0 {
        t.Fatalf("expected exit code 0, got %d\nstdout:\n%s\nstderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
    }
    if !strings.Contains(resp.Stdout, "hello") {
        t.Fatalf("expected stdout to contain 'hello', got:\n%s", resp.Stdout)
    }
}
```
