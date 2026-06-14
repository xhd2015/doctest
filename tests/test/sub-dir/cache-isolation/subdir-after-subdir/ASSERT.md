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
    if !strings.Contains(secondStderr, "1 test case") {
        t.Fatalf("second run expected 1 test case, stderr:\n%s", secondStderr)
    }

    if strings.Contains(secondStderr, "TestGeneratedCaseGroupALeaf1") {
        t.Fatalf("second run expected leaf-1 NOT to run, stderr:\n%s", secondStderr)
    }
    if strings.Contains(secondStderr, "TestGeneratedCaseGroupALeaf2") {
        t.Fatalf("second run expected leaf-2 NOT to run, stderr:\n%s", secondStderr)
    }
    if !strings.Contains(secondStderr, "TestGeneratedCaseGroupBLeaf3") {
        t.Fatalf("second run expected leaf-3 to run, stderr:\n%s", secondStderr)
    }
}
```
