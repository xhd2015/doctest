# Scenario

**Feature**: process env values are not part of the leaf-cache key

```
# leaf Setup may read DOCTEST_CACHE_ENV_PROBE via os.Getenv / LookupEnv / syscall.Getenv
# leaf-cache key is spine-only — env values are NOT mixed in

run A (probe=session-a) -> warm PutPass
run B (probe=session-a) -> Cached hit
run C (probe=session-b) -> still Cached hit (env not in key)
```

## Preconditions
- The go-test-cache root has built the doctest binary and provides tree helpers.
- Each leaf configures how the generated test reads `DOCTEST_CACHE_ENV_PROBE`.
- Product: no osenv value hashing; no go-testlog getenv special-case for Cached.
- Probe var is **not** `DOCTEST_SESSION_ID` so session harness stays stable across A/B.
- **Parallel-safe**: multi-run cfg is a local parameter; results live on
  `req.MRFirst` / `req.MRSecond` (not package globals).

## Steps
1. Build a one-leaf doctest project whose Setup reads `DOCTEST_CACHE_ENV_PROBE`.
2. Warm the leaf-cache with probe value A.
3. Run again with probe value A (expect Cached > 0).
4. Run again with probe value B; still expect Cached > 0 (env not in key).

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

const envCacheProbeVar = "DOCTEST_CACHE_ENV_PROBE"

type envCacheCfg struct {
    TestDir      string
    LeafSetupGo  string
    EnvValueA    string
    EnvValueB    string
}

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    // Clear multi-run fields so a leaked pointer cannot satisfy Assert.
    req.MRFirst = nil
    req.MRSecond = nil
    return nil
}

func createTestTreeWithLeafSetup(dir string, leafSetupGo string) error {
    runCode := "func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) { return &Response{}, nil }"
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

func createEnvProbeTestProject(t *testing.T, dirName string, leafSetupGo string) string {
    t.Helper()
    if leafSetupGo == "" {
        t.Fatal("leafSetupGo is required")
    }
    tmp := t.TempDir()
    if err := os.WriteFile(filepath.Join(tmp, "go.mod"), []byte("module testproj\ngo 1.21\n"), 0644); err != nil {
        t.Fatalf("write go.mod: %v", err)
    }
    testDir := filepath.Join(tmp, dirName)
    if err := os.MkdirAll(testDir, 0755); err != nil {
        t.Fatalf("mkdir test dir: %v", err)
    }
    if err := createTestTreeWithLeafSetup(testDir, leafSetupGo); err != nil {
        t.Fatalf("create test tree: %v", err)
    }
    return testDir
}

func doRunWithEnv(t *testing.T, bin string, args []string, envValue string, extraEnv ...string) *Response {
    t.Helper()
    ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
    defer cancel()

    cmd := exec.CommandContext(ctx, bin, args...)
    cmd.Env = append(os.Environ(), extraEnv...)
    cmd.Env = append(cmd.Env, envCacheProbeVar+"="+envValue)

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

func doEnvCacheRun(t *testing.T, req *Request, cfg envCacheCfg) {
    t.Helper()
    if req.Bin == "" {
        t.Fatal("req.Bin is not set")
    }
    if cfg.EnvValueA == "" {
        cfg.EnvValueA = "session-aaaa-1111"
    }
    if cfg.EnvValueB == "" {
        cfg.EnvValueB = "session-bbbb-2222"
    }

    testDir := cfg.TestDir
    if testDir == "" {
        testDir = createEnvProbeTestProject(t, "envprobe", cfg.LeafSetupGo)
    }

    baseArgs := []string{"test", testDir}

    // Stable GOCACHE + leaf-cache store across warmups and measured runs so
    // only the env probe var differs between A and B.
    goCache := t.TempDir()
    leafCache := t.TempDir()
    stableEnv := []string{
        "GOCACHE=" + goCache,
        "DOCTEST_LEAF_CACHE=" + leafCache,
        "DOCTEST_SESSION_ID=env-probe-stable-session",
    }

    // Two warmups with env A, then measured hit with A, then B (still hit).
    doRunWithEnv(t, req.Bin, baseArgs, cfg.EnvValueA, stableEnv...)
    doRunWithEnv(t, req.Bin, baseArgs, cfg.EnvValueA, stableEnv...)

    req.MRFirst = doRunWithEnv(t, req.Bin, baseArgs, cfg.EnvValueA, stableEnv...)
    req.MRSecond = doRunWithEnv(t, req.Bin, baseArgs, cfg.EnvValueB, stableEnv...)
}
```
