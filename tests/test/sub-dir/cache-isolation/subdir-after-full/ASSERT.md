---
label: heavy
---

## Expected
- First run (full tree) discovers 3 test cases and all 3 pass.
- Second run (group-a sub-dir) discovers 2 test cases.
- Second run stderr contains test function names for leaf-1 and leaf-2.
- Second run stderr does NOT contain test function name for leaf-3 (even though its stale test file remains from the full run).

## Exit Code
- Exit code 0.

```go
import (
    "strings"
    "testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if err != nil {
        t.Fatalf("run failed: %v", err)
    }
    if resp.ExitCode != 0 {
        t.Fatalf("expected exit 0, got %d\nstderr:\n%s", resp.ExitCode, resp.Stderr)
    }
    if req.MRFirst == nil || req.MRSecond == nil {
        t.Fatal("multi-run state not set on req (doMultiRun)")
    }

    firstStderr := req.MRFirst.Stderr
    if !strings.Contains(firstStderr, "3 test cases") {
        t.Fatalf("first run expected 3 test cases, stderr:\n%s", firstStderr)
    }

    secondStdout := req.MRSecond.Stdout
    if !strings.Contains(req.MRSecond.Stderr, "2 test cases") {
        t.Fatalf("second run expected 2 test cases, stderr:\n%s", req.MRSecond.Stderr)
    }

    // Unified: subtests named by path; group-a has leaf-1/leaf-2 only.
    if !strings.Contains(secondStdout, "leaf-1") {
        t.Fatalf("second run expected leaf-1 to run, stdout:\n%s", secondStdout)
    }
    if !strings.Contains(secondStdout, "leaf-2") {
        t.Fatalf("second run expected leaf-2 to run, stdout:\n%s", secondStdout)
    }
    if strings.Contains(secondStdout, "leaf-3") {
        t.Fatalf("second run expected leaf-3 NOT to run, stdout:\n%s", secondStdout)
    }
}
```
