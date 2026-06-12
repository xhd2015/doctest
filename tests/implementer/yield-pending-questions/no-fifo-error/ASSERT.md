## Expected
- Exit code non-zero.
- Stderr contains "QUESTION_FIFO must be set".

```go
import (
    "strings"
    "testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if resp.ExitCode == 0 {
        t.Fatal("expected non-zero exit when QUESTION_FIFO is not set")
    }
    if !strings.Contains(resp.Stderr, "QUESTION_FIFO must be set") {
        t.Fatalf("expected 'QUESTION_FIFO must be set' in stderr:\n%s", resp.Stderr)
    }
}
```
