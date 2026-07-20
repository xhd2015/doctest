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
3. Each leaf sets cfg flags to control behavior (edit path, mtime-only, identical rewrite).

```go
import (
    "bytes"
    "context"
    "errors"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
    "testing"
    "time"

    "github.com/xhd2015/dot-pkgs/go-pkgs/pathfmt"
)

type multiRunCfg struct {
    TestDir           string
    ModifyFile        string
    ModifyContent     string
    UseCountTwo       bool
    TouchMtimeOnly    bool // chtimes ModifyFile only; no content rewrite
    RewriteIdentical  bool // rewrite ModifyFile with its current bytes
}

var cfg multiRunCfg

type multiRunState struct {
    FirstResp  *Response
    SecondResp *Response
    GenDir     string
}

var state multiRunState

func Setup(t *testing.T, req *Request) error {
    // Reset shared package state every leaf — unified suite runs all leaves in
    // one package; stale cfg.ModifyFile from a prior leaf would poison cache tests.
    cfg = multiRunCfg{}
    state = multiRunState{}
    req.Args = []string{}
    return nil
}

func doRun(t *testing.T, bin string, args []string, extraEnv ...string) *Response {
    t.Helper()
    ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
    defer cancel()

    cmd := exec.CommandContext(ctx, bin, args...)
    cmd.Env = append(os.Environ(), extraEnv...)

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
    for _, line := range strings.Split(stderr, "\n") {
        trimmed := strings.TrimSpace(line)
        if !strings.HasPrefix(trimmed, "cd ") || !strings.Contains(trimmed, " && go ") {
            continue
        }
        path := strings.TrimSpace(strings.TrimPrefix(trimmed, "cd "))
        path = strings.Split(path, " && go ")[0]
        path = pathfmt.Expand(path)
        dir := path
        for {
            if strings.Contains(dir, "mapping-gen") {
                if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
                    state.GenDir = dir
                    return
                }
            }
            parent := filepath.Dir(dir)
            if parent == dir {
                break
            }
            dir = parent
        }
        return
    }
}

// stdoutHasPositiveCached reports Cached count > 0 in a summary line.
func stdoutHasPositiveCached(stdout string) bool {
    return strings.Contains(stdout, "Cached") && !strings.Contains(stdout, "0 Cached")
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

    // Isolate GOCACHE + fix session id across warmup/measured runs so concurrent
    // outer-suite noise cannot race the shared machine cache, and nested go test
    // result-cache keys stay stable.
    goCache := t.TempDir()
    stableEnv := []string{
        "GOCACHE=" + goCache,
        "DOCTEST_SESSION_ID=multi-run-stable-session",
    }

    // Two warmups: first generation may rewrite mapping-gen; second stores the
    // go test result for measured runs to hit as cached.
    doRun(t, req.Bin, baseArgs, stableEnv...)
    doRun(t, req.Bin, baseArgs, stableEnv...)

    resp1 := doRun(t, req.Bin, baseArgs, stableEnv...)
    state.FirstResp = resp1

    parseGenDir(resp1.Stderr)

    if cfg.ModifyFile != "" {
        targetPath := filepath.Join(testDir, cfg.ModifyFile)
        switch {
        case cfg.TouchMtimeOnly:
            now := time.Now().Add(2 * time.Second)
            if err := os.Chtimes(targetPath, now, now); err != nil {
                t.Fatalf("chtimes %s: %v", targetPath, err)
            }
        case cfg.RewriteIdentical:
            data, err := os.ReadFile(targetPath)
            if err != nil {
                t.Fatalf("read for identical rewrite: %v", err)
            }
            if err := os.WriteFile(targetPath, data, 0644); err != nil {
                t.Fatalf("identical rewrite: %v", err)
            }
        case cfg.ModifyContent != "":
            if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
                t.Fatalf("mkdir for modify: %v", err)
            }
            if err := os.WriteFile(targetPath, []byte(cfg.ModifyContent), 0644); err != nil {
                t.Fatalf("modify file: %v", err)
            }
        }
    }

    args2 := baseArgs
    if cfg.UseCountTwo {
        args2 = append([]string{"test", testDir, "-count=2"})
    }

    resp2 := doRun(t, req.Bin, args2, stableEnv...)
    state.SecondResp = resp2
}
```
