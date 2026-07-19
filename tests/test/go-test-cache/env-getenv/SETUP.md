# Scenario

**Feature**: prove which environment-variable reads affect Go test result caching

```
# warmup + two captured doctest test runs against a tiny generated leaf
run A (env=session-a) -> cached
run B (env=session-b) -> cache hit or miss depends on how env was read

# os.Getenv / os.LookupEnv -> testlog records getenv -> cache miss on B
# syscall.Getenv       -> no testlog entry      -> cache hit on B
```

## Preconditions
- The go-test-cache root has built the doctest binary and provides tree helpers.
- Each leaf configures how the generated test reads `DOCTEST_SESSION_ID`.

## Steps
1. Build a one-leaf doctest project whose Setup reads `DOCTEST_SESSION_ID`.
2. Warm the go-test cache with env value A.
3. Run again with env value A (expect `1 Cached`).
4. Run again with env value B; cache behavior depends on the read API.

```go
import (
    "bytes"
    "context"
    "errors"
    "os"
    "os/exec"
    "path/filepath"
    "testing"
    "time"

    "github.com/xhd2015/doctest/libdoc/testtree"
)

const envCacheProbeVar = "DOCTEST_SESSION_ID"

type envCacheCfg struct {
    TestDir      string
    LeafSetupGo  string
    EnvValueA    string
    EnvValueB    string
}

var envCfg envCacheCfg

type envCacheState struct {
    FirstResp  *Response
    SecondResp *Response
}

var envState envCacheState

func Setup(t *testing.T, req *Request) error {
    envState.FirstResp = nil
    envState.SecondResp = nil
    if envCfg.EnvValueA == "" {
        envCfg.EnvValueA = "session-aaaa-1111"
    }
    if envCfg.EnvValueB == "" {
        envCfg.EnvValueB = "session-bbbb-2222"
    }
    return nil
}

func createTestTreeWithLeafSetup(dir string, leafSetupGo string) error {
    runCode := "func Run(t *testing.T, req *Request) (*Response, error) { return &Response{}, nil }"
    if err := os.WriteFile(filepath.Join(dir, "DOCTEST.md"), []byte(testtree.MinimalDOCTEST(doctestBody(runCode))), 0644); err != nil {
        return err
    }
    if err := os.WriteFile(filepath.Join(dir, "SETUP.md"), []byte(rootSetupContent()), 0644); err != nil {
        return err
    }
    leafDir := filepath.Join(dir, "simple")
    if err := os.MkdirAll(leafDir, 0755); err != nil {
        return err
    }
    if err := os.WriteFile(filepath.Join(leafDir, "SETUP.md"), []byte(doctestGoBlock(leafSetupGo)), 0644); err != nil {
        return err
    }
    if err := os.WriteFile(filepath.Join(leafDir, "ASSERT.md"), []byte(leafAssertContent()), 0644); err != nil {
        return err
    }
    return nil
}

func createEnvProbeTestProject(t *testing.T, dirName string) string {
    t.Helper()
    if envCfg.LeafSetupGo == "" {
        t.Fatal("envCfg.LeafSetupGo is required")
    }
    tmp := t.TempDir()
    if err := os.WriteFile(filepath.Join(tmp, "go.mod"), []byte("module testproj\ngo 1.21\n"), 0644); err != nil {
        t.Fatalf("write go.mod: %v", err)
    }
    testDir := filepath.Join(tmp, dirName)
    if err := os.MkdirAll(testDir, 0755); err != nil {
        t.Fatalf("mkdir test dir: %v", err)
    }
    if err := createTestTreeWithLeafSetup(testDir, envCfg.LeafSetupGo); err != nil {
        t.Fatalf("create test tree: %v", err)
    }
    return testDir
}

func doRunWithEnv(t *testing.T, bin string, args []string, envValue string) *Response {
    t.Helper()
    ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
    defer cancel()

    cmd := exec.CommandContext(ctx, bin, args...)
    cmd.Env = append(os.Environ(), envCacheProbeVar+"="+envValue)

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

func doEnvCacheRun(t *testing.T, req *Request) {
    t.Helper()
    if req.Bin == "" {
        t.Fatal("req.Bin is not set")
    }

    testDir := envCfg.TestDir
    if testDir == "" {
        testDir = createEnvProbeTestProject(t, "envprobe")
    }

    baseArgs := []string{"test", testDir}

    // Two warmups with env A: first gen rewrite may use -count=1; second stores
    // a cache entry for the "1 Cached" first measured run.
    doRunWithEnv(t, req.Bin, baseArgs, envCfg.EnvValueA)
    doRunWithEnv(t, req.Bin, baseArgs, envCfg.EnvValueA)

    envState.FirstResp = doRunWithEnv(t, req.Bin, baseArgs, envCfg.EnvValueA)
    envState.SecondResp = doRunWithEnv(t, req.Bin, baseArgs, envCfg.EnvValueB)
}
```