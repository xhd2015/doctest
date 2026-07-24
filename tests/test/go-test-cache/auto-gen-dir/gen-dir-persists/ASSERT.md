---
label: heavy
---

## Expected
- The gen directory exists and is not empty after the first run completes.
- The gen directory path was captured on `req.MRGenDir`.

## Exit Code
- Exit code 0.

```go
import (
    "os"
    "strings"
    "testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
    if req.MRGenDir == "" {
        t.Fatal("gen dir path was not captured from output")
    }
    if !strings.Contains(req.MRGenDir, "/doctest/") {
        t.Fatalf("gen dir is not a hash-based doctest path; got: %s", req.MRGenDir)
    }
    fi, statErr := os.Stat(req.MRGenDir)
    if statErr != nil {
        t.Fatalf("gen dir does not exist at %s: %v", req.MRGenDir, statErr)
    }
    if !fi.IsDir() {
        t.Fatalf("gen dir path is not a directory: %s", req.MRGenDir)
    }
    entries, readErr := os.ReadDir(req.MRGenDir)
    if readErr != nil {
        t.Fatalf("cannot read gen dir: %v", readErr)
    }
    if len(entries) == 0 {
        t.Fatal("gen dir is empty")
    }
}
```
