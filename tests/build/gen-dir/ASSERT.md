## Expected
- The command succeeds.
- Generated Go files are written to the requested directory.

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
    for _, name := range []string{"happy_path_test.go", "expected_error_test.go", "override_run_test.go"} {
        if _, err := os.Stat(filepath.Join(genDir, name)); err != nil {
            t.Fatalf("generated file %s missing: %v\nstderr:\n%s", name, err, resp.Stderr)
        }
    }
}
```
