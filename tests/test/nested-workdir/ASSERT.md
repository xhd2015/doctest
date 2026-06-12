## Expected
- The command succeeds from the nested working directory.
- The wrapper resolves the underlying runner independent of the caller's cwd.

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
        t.Fatalf("expected zero exit, got %d\nstdout:\n%s\nstderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
    }
    if strings.Contains(resp.Stderr, "agents/doctest/agents/doctest") {
        t.Fatalf("runner path was resolved relative to nested cwd:\n%s", resp.Stderr)
    }
    if strings.Contains(resp.Stderr, "directory not found") {
        t.Fatalf("runner failed with directory-not-found from nested cwd:\n%s", resp.Stderr)
    }
}
```
