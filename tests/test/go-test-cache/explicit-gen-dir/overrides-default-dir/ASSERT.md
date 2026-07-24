---
label: heavy
---

## Expected
- The command succeeds.
- The gen dir path printed in stderr matches the user-provided --gen-dir.
- The gen dir is NOT under the default hash-based path.

## Exit Code
- Exit code 0.

```go
import (
    "os"
    "strings"
    "testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
    if err != nil {
        t.Fatalf("run failed: %v", err)
    }
    if resp.ExitCode != 0 {
        t.Fatalf("expected exit 0, got %d\nstderr:\n%s", resp.ExitCode, resp.Stderr)
    }
    genDir := ""
    for _, env := range req.Env {
        if len(env) > 8 && env[:8] == "GEN_DIR=" {
            genDir = env[8:]
            break
        }
    }
    if genDir == "" {
        t.Fatal("GEN_DIR not set in env")
    }

    genPrefix := "→ " + genDir
    if !strings.Contains(resp.Stderr, genPrefix) {
        t.Fatalf("expected stderr to contain gen dir %q, got:\n%s", genPrefix, resp.Stderr)
    }

    cacheDir, cacheErr := os.UserCacheDir()
    if cacheErr == nil {
        doctestCacheDir := cacheDir + "/doctest/"
        if strings.Contains(resp.Stderr, doctestCacheDir) {
            t.Fatalf("gen dir path contains default hash-based cache dir %q when --gen-dir was specified:\n%s", doctestCacheDir, resp.Stderr)
        }
    }
}
```
