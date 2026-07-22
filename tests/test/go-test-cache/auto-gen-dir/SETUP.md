# Scenario

**Feature**: multi-run harness for leaf-cache Cached product under auto-gen-dir

```
# warmup stores leaf-cache pass markers (+ may warm go package cache)
doctest test <fixture> -> leaf pass -> PutPass(key)

# measured run 1: spine unchanged -> leaf-cache skip -> Cached > 0
# measured run 2: optional spine edit / -count / -a -> hit or miss per product
doctest test <fixture> [flags] -> summary (N Run, N Pass, N Fail, N Cached)
```

## Preconditions
- The root Setup has built the doctest binary and set req.Bin.
- The go-test-cache root has set a 120s timeout.
- Warmup runs populate leaf-cache (and go cache) before the two captured runs.
- **Cached** means leaf-cache skips when the package runs; whole-package go
  `(cached)` expands to N Cached for all N leaves.
- **Parallel-safe**: multi-run cfg is a local parameter; results live on
  `req.MRFirst` / `req.MRSecond` / `req.MRGenDir` (not package globals).

## Steps
1. Provide a multi-run helper that takes cfg, orchestrates silent warmups, then
   two captured executions, and writes results onto `req`.
2. Each leaf builds a local `multiRunCfg` (edit path, mtime-only, identical rewrite,
   SecondFlags) and calls `doMultiRun(t, req, cfg)`.
3. Use a unique `DOCTEST_SESSION_ID` per multi-run instance (stable across that
   instance's warmup/measured pair; isolated from concurrent leaves).

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
    // SecondFlags are appended to `doctest test <dir>` on the second measured run
    // (e.g. "-count=1", "-a"). Any -count or -a must bypass leaf-cache → 0 Cached.
    SecondFlags       []string
    TouchMtimeOnly    bool // chtimes ModifyFile only; no content rewrite
    RewriteIdentical  bool // rewrite ModifyFile with its current bytes
}

func Setup(t *testing.T, req *Request) error {
    // Clear multi-run fields so a leaked pointer from a prior leaf cannot
    // satisfy Assert under sequential reuse of a Request value (defensive;
    // doctest allocates per leaf — package globals were the real race).
    req.MRFirst = nil
    req.MRSecond = nil
    req.MRGenDir = ""
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

func parseGenDir(stderr string) string {
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
                    return dir
                }
            }
            parent := filepath.Dir(dir)
            if parent == dir {
                break
            }
            dir = parent
        }
        return ""
    }
    return ""
}

// stdoutHasPositiveCached is defined on the go-test-cache root SETUP.

func doMultiRun(t *testing.T, req *Request, cfg multiRunCfg) {
    if req.Bin == "" {
        t.Fatalf("req.Bin is not set")
    }

    testDir := cfg.TestDir
    if testDir == "" {
        testDir = createTempTestProject(t, "mytest")
    }

    baseArgs := []string{"test", testDir}

    // Isolate GOCACHE + leaf-cache store; unique session id per multi-run so
    // concurrent leaves do not share session/leaf-cache state.
    goCache := t.TempDir()
    leafCache := t.TempDir()
    sessionID := "multi-run-" + filepath.Base(leafCache)
    stableEnv := []string{
        "GOCACHE=" + goCache,
        "DOCTEST_LEAF_CACHE=" + leafCache,
        "DOCTEST_SESSION_ID=" + sessionID,
    }

    // Two warmups: first generation may rewrite mapping-gen; second stores the
    // go test result for measured runs to hit as cached.
    doRun(t, req.Bin, baseArgs, stableEnv...)
    doRun(t, req.Bin, baseArgs, stableEnv...)

    resp1 := doRun(t, req.Bin, baseArgs, stableEnv...)
    req.MRFirst = resp1
    req.MRGenDir = parseGenDir(resp1.Stderr)

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
    if len(cfg.SecondFlags) > 0 {
        args2 = append([]string{"test", testDir}, cfg.SecondFlags...)
    }

    resp2 := doRun(t, req.Bin, args2, stableEnv...)
    req.MRSecond = resp2
}
```
