## Expected
- Exit code is 0 (test compiled and passed).
- The unused `"fmt"` import was removed by `imports.Process` during code generation.
- No compilation errors in the output.

```go
import (
    "testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if err != nil {
        t.Fatalf("run error: %v", err)
    }

    if resp.ExitCode != 0 {
        t.Fatalf("expected exit code 0 (compiled and passed), got %d\nstdout: %s\nstderr: %s", resp.ExitCode, resp.Stdout, resp.Stderr)
    }
}
```
