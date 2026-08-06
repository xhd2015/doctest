## Expected
- First run (group-a sub-dir) discovers 2 test cases and both pass.
- Second run (full tree) discovers 3 test cases.
- Second run stderr contains test function names for all 3 leaves (leaf-1, leaf-2, leaf-3).
- This verifies that a sub-dir cache does not prevent a subsequent full run from executing all tests.

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
    if !strings.Contains(req.MRSecond.Stderr, "3 test cases") {
        t.Fatalf("second run expected 3 test cases, stderr:\n%s", req.MRSecond.Stderr)
    }

    // Unified suite subtest names use leaf paths (slashes encoded as __).
    for _, leaf := range []string{"leaf-1", "leaf-2", "leaf-3"} {
        if !strings.Contains(secondStdout, leaf) {
            t.Fatalf("second run expected %s to appear in stdout:\n%s", leaf, secondStdout)
        }
    }
}
```
