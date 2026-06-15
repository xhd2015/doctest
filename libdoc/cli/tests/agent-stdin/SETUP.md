## Preconditions
- The doctest binary is built and available.
- Tests run the CLI directly by calling `cli.Run()` after replacing `os.Stdin`.

## Steps
1. Child SETUP.md files set up stdin source, args, and requirement files.
2. Run replaces `os.Stdin` (pipe, /dev/null, or keep original), calls `cli.Run(req.Args)`, captures stdout/stderr/error.

```go
import (
    "bytes"
    "io"
    "os"
    "testing"

    "github.com/xhd2015/doctest/libdoc/cli"
)

type Request struct {
    Args        []string
    Stdin       string
    StdinSource string
    ReqContent  string
}

type Response struct {
    Err    error
    Stdout string
    Stderr string
}

func Run(t *testing.T, req *Request) (*Response, error) {
    oldStdin := os.Stdin
    defer func() { os.Stdin = oldStdin }()

    switch req.StdinSource {
    case "pipe":
        r, w, err := os.Pipe()
        if err != nil {
            return nil, err
        }
        if req.Stdin != "" {
            if _, err := w.WriteString(req.Stdin); err != nil {
                w.Close()
                r.Close()
                return nil, err
            }
        }
        w.Close()
        os.Stdin = r
        defer r.Close()
    case "devnull":
        f, err := os.Open(os.DevNull)
        if err != nil {
            return nil, err
        }
        os.Stdin = f
        defer f.Close()
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

    return &Response{
        Err:    cliErr,
        Stdout: stdoutBuf.String(),
        Stderr: stderrBuf.String(),
    }, nil
}
```
