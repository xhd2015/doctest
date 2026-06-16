# Scenario

**Feature**: the root Setup has built the doctest binary and set req.Bin

```
# build and run test binary, report results
doctest test <dir> -> build -> run binary -> pass/fail per leaf -> exit code

# path patterns
.../ -> walk tree | subdir -> run subtree | multi-dir -> aggregate results

# output
progress dots -> . F | verbose -> go test -v | count -> N tests
```

## Preconditions
- The root Setup has built the doctest binary and set req.Bin.
- The go-test-cache root has set a 120s timeout.
- A warmup run populates the cache before the two captured runs.

## Steps
1. Define shared configuration and state for multi-run test leaves.
2. Provide a multi-run `Run` function that orchestrates a silent warmup run followed by two captured executions.
3. Each leaf sets cfg flags to control behavior.

```go
import (
    "bytes"
    "context"
    "errors"
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
    "testing"
    "time"
)

type multiRunCfg struct {
    TestDir        string
    ModifyFile     string
    ModifyContent  string
    UseCountTwo    bool
}

var cfg multiRunCfg

type multiRunState struct {
    FirstResp  *Response
    SecondResp *Response
    GenDir     string
}

var state multiRunState

func Setup(t *testing.T, req *Request) error {
    state.GenDir = ""
    state.FirstResp = nil
    state.SecondResp = nil
    req.Args = []string{}
    return nil
}

func doRun(t *testing.T, bin string, args []string) *Response {
    t.Helper()
    ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
    defer cancel()

    cmd := exec.CommandContext(ctx, bin, args...)
    cmd.Env = os.Environ()

    var stdoutBuf bytes.Buffer
    var stderrBuf bytes.Buffer
    cmd.Stdout = &stdoutBuf
    cmd.Stderr = &stderrBuf

    resp := &Response{}
    err := cmd.Run()
    resp.Stdout = stdoutBuf.String()
    resp.Stderr = stderrBuf.String()
    if err == nil {
        return resp
    }
    var exitErr *exec.ExitError
    if errors.As(err, &exitErr) {
        resp.ExitCode = exitErr.ExitCode()
        return resp
    }
    resp.Err = err
    return resp
}

func parseGenDir(stderr string) {
    idx := strings.Index(stderr, "→ ")
    if idx < 0 {
        return
    }
    rest := stderr[idx+len("→ "):]
    end := strings.IndexFunc(rest, func(r rune) bool { return r == '\n' || r == ' ' })
    if end < 0 {
        state.GenDir = strings.TrimSpace(rest)
    } else {
        state.GenDir = strings.TrimSpace(rest[:end])
    }
    _ = fmt.Sprintf
}

func doMultiRun(t *testing.T, req *Request) {
    if req.Bin == "" {
        t.Fatalf("req.Bin is not set")
    }

    testDir := cfg.TestDir
    if testDir == "" {
        testDir = createTempTestProject(t, "mytest")
    }

    baseArgs := []string{"test", testDir}

    doRun(t, req.Bin, baseArgs)

    resp1 := doRun(t, req.Bin, baseArgs)
    state.FirstResp = resp1

    parseGenDir(resp1.Stderr)

    if cfg.ModifyFile != "" && cfg.ModifyContent != "" {
        targetPath := filepath.Join(testDir, cfg.ModifyFile)
        if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
            t.Fatalf("mkdir for modify: %v", err)
        }
        if err := os.WriteFile(targetPath, []byte(cfg.ModifyContent), 0644); err != nil {
            t.Fatalf("modify file: %v", err)
        }
    }

    args2 := baseArgs
    if cfg.UseCountTwo {
        args2 = append([]string{"test", testDir, "-count=2"})
    }

    resp2 := doRun(t, req.Bin, args2)
    state.SecondResp = resp2
}
```
