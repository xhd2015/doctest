---
label: heavy
---

## Expected
- The gen directory exists and is not empty after the first run completes.
- The gen directory path was captured in state.

## Exit Code
- Exit code 0.

```go
import (
    "os"
    "strings"
    "testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if state.GenDir == "" {
        t.Fatal("gen dir path was not captured from output")
    }
    if !strings.Contains(state.GenDir, "/doctest/") {
        t.Fatalf("gen dir is not a hash-based doctest path; got: %s", state.GenDir)
    }
    fi, statErr := os.Stat(state.GenDir)
    if statErr != nil {
        t.Fatalf("gen dir does not exist at %s: %v", state.GenDir, statErr)
    }
    if !fi.IsDir() {
        t.Fatalf("gen dir path is not a directory: %s", state.GenDir)
    }
    entries, readErr := os.ReadDir(state.GenDir)
    if readErr != nil {
        t.Fatalf("cannot read gen dir: %v", readErr)
    }
    if len(entries) == 0 {
        t.Fatal("gen dir is empty")
    }
}
```
