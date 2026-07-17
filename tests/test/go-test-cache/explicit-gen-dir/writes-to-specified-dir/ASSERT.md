---
label: heavy
---

## Expected
- The command succeeds.
- Generated Go test files exist at the specified --gen-dir path.

## Exit Code
- Exit code 0.

```go
import (
    "os"
    "path/filepath"
    "testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
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
    foundTestFile := false
    walkErr := filepath.WalkDir(genDir, func(path string, d os.DirEntry, err error) error {
        if err != nil {
            return err
        }
        if filepath.Ext(d.Name()) == ".go" {
            foundTestFile = true
            return filepath.SkipAll
        }
        return nil
    })
    if walkErr != nil {
        t.Fatalf("cannot walk gen dir %s: %v\nstderr:\n%s", genDir, walkErr, resp.Stderr)
    }
    if !foundTestFile {
        t.Fatalf("no .go files found in gen dir %s\nstderr:\n%s", genDir, resp.Stderr)
    }
}
```
