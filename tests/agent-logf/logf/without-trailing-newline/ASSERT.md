## Expected
- Output has `[YYYY-MM-DDTHH:MM:SS] hello` followed by exactly one newline.
- No double newline.

```go
import (
    "regexp"
    "strings"
    "testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if err != nil {
        t.Fatalf("run failed: %v", err)
    }

    tsPattern := regexp.MustCompile(`^\[\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\] hello\n$`)
    if !tsPattern.MatchString(resp.Stdout) {
        t.Fatalf("expected timestamped 'hello' with one newline, got:\n%q", resp.Stdout)
    }

    if strings.Count(resp.Stdout, "\n") != 1 {
        t.Fatalf("expected exactly one newline, got %d in:\n%q", strings.Count(resp.Stdout, "\n"), resp.Stdout)
    }
}
```
