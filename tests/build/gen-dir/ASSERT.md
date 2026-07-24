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

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
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
    // Unified layout: non-test leaf.go per leaf + suite_test.go + droot.go.
    leafGo := 0
    hasSuite := false
    hasDroot := false
    for _, f := range found {
        switch f {
        case "leaf.go":
            leafGo++
        case "suite_test.go":
            hasSuite = true
        case "droot.go":
            hasDroot = true
        }
    }
    if leafGo < 3 {
        t.Fatalf("expected ≥3 leaf.go files, found %d; all: %v\nstderr:\n%s", leafGo, found, resp.Stderr)
    }
    if !hasSuite {
        t.Fatalf("expected suite_test.go; found: %v\nstderr:\n%s", found, resp.Stderr)
    }
    if !hasDroot {
        t.Fatalf("expected droot.go; found: %v\nstderr:\n%s", found, resp.Stderr)
    }
}
```
