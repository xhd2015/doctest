## Preconditions
- A doctest binary is available at `req.Bin` (built by the root SETUP.md).
- Session directories are created under a temp `DOCTEST_DEBUG_SESSION_HOME`.

## Steps
1. Run the doctest binary with args set by the leaf (e.g., `agent implement --status --session-id X`).
2. Capture stdout, stderr, and exit code.

```go
import (
    "bytes"
    "context"
    "errors"
    "fmt"
    "os/exec"
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    req.Env = append(req.Env, "TEST_GROUP=show-status")
    return nil
}

func Run(t *testing.T, req *Request) (*Response, error) {
    ctx, cancel := context.WithTimeout(context.Background(), req.Timeout)
    defer cancel()

    cmd := exec.CommandContext(ctx, req.Bin, req.Args...)
    var stdout, stderr bytes.Buffer
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr
    err := cmd.Run()

    resp := &Response{
        Stdout: stdout.String(),
        Stderr: stderr.String(),
        Err:    err,
    }
    if err == nil {
        return resp, nil
    }
    var exitErr *exec.ExitError
    if errors.As(err, &exitErr) {
        resp.ExitCode = exitErr.ExitCode()
        return resp, nil
    }
    if ctx.Err() != nil {
        return resp, ctx.Err()
    }
    return resp, fmt.Errorf("command failed: %w", err)
}
```
