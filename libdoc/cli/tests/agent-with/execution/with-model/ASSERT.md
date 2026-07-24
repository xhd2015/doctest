## Expected
- Exit code 0.
- Stdout contains `gpt-4`.

## Side Effects
- The child process sees `DOCTEST_SUBAGENT_MODEL=gpt-4` and `DOCTEST_SUBAGENT_AGENT_RUNNER=opencode`.

## Exit Code
- 0.

```go
import (
    "strings"
    "testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
    if resp.Err != nil {
        t.Fatalf("unexpected error: %v", resp.Err)
    }
    if resp.ExitCode != 0 {
        t.Fatalf("expected exit code 0, got %d\nstdout:\n%s\nstderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
    }
    stdout := strings.TrimSpace(resp.Stdout)
    if stdout != "gpt-4" {
        t.Fatalf("expected stdout 'gpt-4', got: %q", stdout)
    }
}
```
