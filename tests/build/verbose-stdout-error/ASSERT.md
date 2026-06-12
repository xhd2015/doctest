## Expected
- With the fix applied, `doctest test -v` runs successfully (exit 0)
  because the verbose code path now uses `Run()` instead of `CombinedOutput()`
  when Stdout/Stderr are already set.

```go
import (
    "strings"
    "testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if err != nil {
        t.Fatalf("run failed unexpectedly: %v", err)
    }
    if resp.ExitCode != 0 {
        t.Fatalf("expected doctest test -v to succeed (exit 0)\nstderr:\n%s", resp.Stderr)
    }
}
```
