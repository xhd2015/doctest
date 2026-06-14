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
    files, readErr := os.ReadDir(genDir)
    if readErr != nil {
        t.Fatalf("cannot read gen dir %s: %v\nstderr:\n%s", genDir, readErr, resp.Stderr)
    }
    if len(files) == 0 {
        t.Fatalf("gen dir is empty, expected generated test files\nstderr:\n%s", resp.Stderr)
    }
    foundTestFile := false
    for _, f := range files {
        if filepath.Ext(f.Name()) == ".go" {
            foundTestFile = true
            break
        }
    }
    if !foundTestFile {
        t.Fatalf("no .go files found in gen dir %s\nfiles: %v\nstderr:\n%s", genDir, files, resp.Stderr)
    }
}
```
