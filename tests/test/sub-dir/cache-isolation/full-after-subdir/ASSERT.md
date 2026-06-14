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

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if err != nil {
        t.Fatalf("run failed: %v", err)
    }
    if resp.ExitCode != 0 {
        t.Fatalf("expected exit 0, got %d\nstderr:\n%s", resp.ExitCode, resp.Stderr)
    }
    if isoState.FirstRun == nil || isoState.SecondRun == nil {
        t.Fatal("multi-run state not set")
    }

    firstStderr := isoState.FirstRun.Stderr
    if !strings.Contains(firstStderr, "2 test cases") {
        t.Fatalf("first run expected 2 test cases, stderr:\n%s", firstStderr)
    }

    secondStderr := isoState.SecondRun.Stderr
    if !strings.Contains(secondStderr, "3 test cases") {
        t.Fatalf("second run expected 3 test cases, stderr:\n%s", secondStderr)
    }

    if !strings.Contains(secondStderr, "TestGeneratedCaseGroupALeaf1") {
        t.Fatalf("second run expected leaf-1 to run, stderr:\n%s", secondStderr)
    }
    if !strings.Contains(secondStderr, "TestGeneratedCaseGroupALeaf2") {
        t.Fatalf("second run expected leaf-2 to run, stderr:\n%s", secondStderr)
    }
    if !strings.Contains(secondStderr, "TestGeneratedCaseGroupBLeaf3") {
        t.Fatalf("second run expected leaf-3 to run, stderr:\n%s", secondStderr)
    }
}
```
