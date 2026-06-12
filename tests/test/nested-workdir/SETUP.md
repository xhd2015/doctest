## Preconditions
- A valid doc-style test tree exists in the module.
- The command is invoked from inside a nested module directory.

## Steps
1. Set the process working directory to the module root (`DOCTEST_ROOT/..`).
2. Run `doctest test <absolute-dir>`.

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
    exampleDir := filepath.Join(DOCTEST_ROOT, "testdata/basic-request-runner")
    req.WorkDir = filepath.Join(DOCTEST_ROOT, "..")
    req.Args = []string{"test", exampleDir}

    tmp := t.TempDir()
    doctestBin := filepath.Join(tmp, "doctest")
    buildDT := exec.Command("go", "build", "-o", doctestBin, "./cmd/doctest")
    buildDT.Dir = filepath.Join(DOCTEST_ROOT, "..")
    if out, err := buildDT.CombinedOutput(); err != nil {
        t.Fatalf("build doctest: %v\n%s", err, string(out))
    }
    req.Bin = doctestBin
    return nil
}

func Run(t *testing.T, req *Request) (*Response, error) {
    ctx, cancel := context.WithTimeout(context.Background(), req.Timeout)
    defer cancel()

    bin := req.Bin
    if bin == "" {
        return nil, fmt.Errorf("req.Bin is not set")
    }
    cmd := exec.CommandContext(ctx, bin, req.Args...)
    cmd.Dir = req.WorkDir
    cmd.Env = append(os.Environ(), req.Env...)

    var stdout bytes.Buffer
    var stderr bytes.Buffer
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr

    err := cmd.Run()
    resp := &Response{
        Stdout: stdout.String(),
        Stderr: stderr.String(),
        Err: err,
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
