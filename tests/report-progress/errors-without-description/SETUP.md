# Scenario

**Feature**: `PROGRESS_FILE` is set but no description argument is provided

```
# sub-agents report progress to a file
sub-agent --writes--> progress file (env var DOCTEST_PROGRESS_FILE)

# multiple entries append
each step -> structured JSON entry -> append to file
```

## Preconditions
- `PROGRESS_FILE` is set but no description argument is provided.

## Steps
1. Set `PROGRESS_FILE` to a temp file.
2. Invoke `report-progress` with no arguments.

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
    req.Args = nil // no description
    return nil
}
```
