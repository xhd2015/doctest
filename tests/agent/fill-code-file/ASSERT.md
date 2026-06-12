## Expected
- The command fails because the target is not a directory.

```go
import (
    "strings"
    "testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if err != nil {
        t.Fatalf("run failed unexpectedly: %v", err)
    }
    if resp.ExitCode == 0 {
        t.Fatalf("expected nonzero exit")
    }
    if !strings.Contains(resp.Stderr, "not a directory") {
        t.Fatalf("stderr missing not-a-directory error:\n%s", resp.Stderr)
    }
}
```
