---
label: heavy
---

## Expected
- First run (group-a sub-dir) discovers 2 test cases and both pass.
- Second run (group-b sub-dir) discovers 1 test case.
- Second run stderr contains test function name for leaf-3.
- Second run stderr does NOT contain test function names for leaf-1 or leaf-2.

## Exit Code
- Exit code 0.

```go
import (
    "strings"
    "testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
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
    if !strings.Contains(firstStderr, "2 test cases") {
        t.Fatalf("first run expected 2 test cases, stderr:\n%s", firstStderr)
    }

    secondStdout := req.MRSecond.Stdout
    if !strings.Contains(req.MRSecond.Stderr, "1 test case") {
        t.Fatalf("second run expected 1 test case, stderr:\n%s", req.MRSecond.Stderr)
    }

    if strings.Contains(secondStdout, "leaf-1") {
        t.Fatalf("second run expected leaf-1 NOT to run, stdout:\n%s", secondStdout)
    }
    if strings.Contains(secondStdout, "leaf-2") {
        t.Fatalf("second run expected leaf-2 NOT to run, stdout:\n%s", secondStdout)
    }
    if !strings.Contains(secondStdout, "leaf-3") {
        t.Fatalf("second run expected leaf-3 to run, stdout:\n%s", secondStdout)
    }
}
```
