## Preconditions
- The auto-detection logic is implemented in `github.com/xhd2015/agent-pro/agent/subagent`.
- The CLI defaults in `cli.go` are changed to empty `Options{}`.
- Tests call the CLI directly via `cli.Run()` and capture the returned error.

## Steps
1. Child SETUP.md files configure Args and Env for each test scenario.
2. Run sets env vars via os.Setenv, calls `cli.Run(req.Args)`, captures stdout/stderr and error.

```go
import (
    "bytes"
    "errors"
    "io"
    "os"
    "os/exec"
    "strings"
    "testing"

    "github.com/xhd2015/doctest/libdoc/cli"
)

type Request struct {
    Args []string
    Env  []string
}

type Response struct {
    ExitCode int
    Stdout   string
    Stderr   string
    Err      error
}

func Run(t *testing.T, req *Request) (*Response, error) {
    for _, e := range req.Env {
        parts := strings.SplitN(e, "=", 2)
        if len(parts) == 2 {
            os.Setenv(parts[0], parts[1])
        }
    }

    oldStdout := os.Stdout
    rOut, wOut, err := os.Pipe()
    if err != nil {
        return nil, err
    }
    os.Stdout = wOut

    oldStderr := os.Stderr
    rErr, wErr, err := os.Pipe()
    if err != nil {
        wOut.Close()
        rOut.Close()
        return nil, err
    }
    os.Stderr = wErr

    cliErr := cli.Run(req.Args)

    wOut.Close()
    wErr.Close()

    var stdoutBuf, stderrBuf bytes.Buffer
    io.Copy(&stdoutBuf, rOut)
    io.Copy(&stderrBuf, rErr)
    rOut.Close()
    rErr.Close()

    os.Stdout = oldStdout
    os.Stderr = oldStderr

    resp := &Response{
        Stdout: stdoutBuf.String(),
        Stderr: stderrBuf.String(),
        Err:    cliErr,
    }
    if cliErr == nil {
        return resp, nil
    }
    var exitErr *exec.ExitError
    if errors.As(cliErr, &exitErr) {
        resp.ExitCode = exitErr.ExitCode()
        return resp, nil
    }
    return resp, nil
}
```
