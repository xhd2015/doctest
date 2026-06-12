## Expected
- The command succeeds.
- stdout or stderr proves the fake Codex response was used.

## Side Effects
- The target directory exists after generation.

## Exit Code
- Exit code 0.

```go
import (
    "os"
    "strings"
    "testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if err != nil {
        t.Fatalf("run failed: %v", err)
    }
    if resp.ExitCode != 0 {
        t.Fatalf("exit code = %d\nstdout:\n%s\nstderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
    }
    outDir := req.Args[len(req.Args)-1]
    if _, statErr := os.Stat(outDir); statErr != nil {
        t.Fatalf("generated dir missing: %v", statErr)
    }
    combined := resp.Stdout + "\n" + resp.Stderr
    if !strings.Contains(combined, "fake doctest agent completed") {
        t.Fatalf("output does not show fake codex response:\n%s", combined)
    }
}
```
