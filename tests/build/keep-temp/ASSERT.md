## Expected
- Exit 0.
- The temp directory path (first line of stderr) still exists after the command exits.

```go
import (
    "os"
    "strings"
    "testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if err != nil {
        t.Fatalf("run failed unexpectedly: %v", err)
    }
    if resp.ExitCode != 0 {
        t.Fatalf("expected zero exit, got %d\nstdout:\n%s\nstderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
    }

    firstLine, _, _ := strings.Cut(resp.Stderr, "\n")
    tmpDir := strings.TrimPrefix(firstLine, "→ ")
    if tmpDir == "" || tmpDir == firstLine {
        t.Fatalf("could not parse temp dir from stderr:\n%s", resp.Stderr)
    }

    if _, err := os.Stat(tmpDir); err != nil {
        t.Fatalf("temp dir should persist without --rm, but was deleted: %s\nstderr:\n%s", tmpDir, resp.Stderr)
    }
}
```
