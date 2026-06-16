# Scenario

**Feature**: tests in this group run the `report-progress` binary directly

```
# sub-agents report progress to a file
sub-agent --writes--> progress file (env var DOCTEST_PROGRESS_FILE)

# multiple entries append
each step -> structured JSON entry -> append to file
```

## Preconditions
- Tests in this group run the `report-progress` binary directly.
- The binary is dispatched via the doctest binary copied as `report-progress`.

## Steps
1. Copy the built doctest binary to a temp dir as `report-progress`.
2. Run the binary with the leaf's args and env.

```go
import (
    "bytes"
    "context"
    "errors"
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    req.Env = append(req.Env, "TEST_GROUP=report-progress")
    return nil
}

func Run(t *testing.T, req *Request) (*Response, error) {
    tmp := t.TempDir()
    rpBin := filepath.Join(tmp, "report-progress")
    if out, err := exec.Command("cp", req.Bin, rpBin).CombinedOutput(); err != nil {
        return nil, fmt.Errorf("copy report-progress: %w\n%s", err, string(out))
    }

    ctx, cancel := context.WithTimeout(context.Background(), req.Timeout)
    defer cancel()

    cmd := exec.CommandContext(ctx, rpBin, req.Args...)
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
