---
label: heavy
---

## Expected
- The command succeeds.
- Generated Go files are written to the requested directory in per-leaf structure.

```go
import (
    "os"
    "path/filepath"
    "testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if err != nil {
        t.Fatalf("run failed unexpectedly: %v", err)
    }
    if resp.ExitCode != 0 {
        t.Fatalf("expected zero exit, got %d\nstdout:\n%s\nstderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
    }
    genDir := req.Args[len(req.Args)-1]
    var found []string
    filepath.Walk(genDir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if !info.IsDir() && filepath.Ext(path) == ".go" {
            found = append(found, filepath.Base(path))
        }
        return nil
    })
    for _, name := range []string{"happy_path_test.go", "expected_error_test.go", "override_run_test.go"} {
        ok := false
        for _, f := range found {
            if f == name {
                ok = true
                break
            }
        }
        if !ok {
            t.Fatalf("generated file %s missing; found: %v\nstderr:\n%s", name, found, resp.Stderr)
        }
    }
}
```
