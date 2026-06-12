## Preconditions
- Tests in this group run the `yield-pending-questions` binary directly.
- The binary is dispatched via the doctest binary copied as `yield-pending-questions`.

## Steps
1. Read `YIELD_PQ_BIN` from environment.
2. Run the binary with the leaf's args.

```go
import (
    "bytes"
    "context"
    "errors"
    "fmt"
    "os"
    "os/exec"
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    req.Env = append(req.Env, "TEST_GROUP=yield-pending-questions")
    return nil
}

func Run(t *testing.T, req *Request) (*Response, error) {
    bin := os.Getenv("YIELD_PQ_BIN")
    if bin == "" {
        return nil, fmt.Errorf("YIELD_PQ_BIN not set")
    }

    ctx, cancel := context.WithTimeout(context.Background(), req.Timeout)
    defer cancel()

    cmd := exec.CommandContext(ctx, bin, req.Args...)
    cmd.Env = append(os.Environ(), req.Env...)

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
    return resp, err
}
```
