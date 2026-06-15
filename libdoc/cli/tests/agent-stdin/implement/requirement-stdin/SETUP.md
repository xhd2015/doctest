## Preconditions
- `--requirement` flag with file content and stdin data (no positional args).

## Steps
1. Write a temp requirement file, append `--requirement` flag.
2. StdinSource = "pipe", Stdin = "additional context".

```go
import (
    "os"
    "path/filepath"
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    tf := filepath.Join(t.TempDir(), "requirement.md")
    if err := os.WriteFile(tf, []byte("Feature spec from requirement file"), 0644); err != nil {
        return err
    }
    req.Args = append(req.Args, "--requirement", tf)
    req.StdinSource = "pipe"
    req.Stdin = "additional context from stdin"
    return nil
}
```
