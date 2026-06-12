## Preconditions
- `PROGRESS_FILE` env var points to a writable file.
- The binary receives a description argument.

## Steps
1. Set `PROGRESS_FILE` to a temp file path.
2. Invoke `report-progress` with a description string.
3. Read the file and verify a progress entry was written.

```go
import (
    "os"
    "path/filepath"
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    tf := filepath.Join(t.TempDir(), "progress.jsonl")
    req.Env = append(req.Env, "PROGRESS_FILE="+tf)
    os.Setenv("TEST_PROGRESS_FILE", tf)
    req.Args = []string{"implementing JSON parser"}
    return nil
}
```
