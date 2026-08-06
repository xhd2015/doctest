## Expected
- Output has `[YYYY-MM-DDTHH:MM:SS] item=alpha count=42` followed by exactly one newline.
- No double newline appended (message already ends with `\n`).
- Format verbs `%s` are correctly resolved with the string arguments.

```go
import (
    "regexp"
    "strings"
    "testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
    if err != nil {
        t.Fatalf("run failed: %v", err)
    }

    tsPattern := regexp.MustCompile(`^\[\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\] item=alpha count=42\n$`)
    if !tsPattern.MatchString(resp.Stdout) {
        t.Fatalf("expected timestamped 'item=alpha count=42' with exactly one newline, got:\n%q", resp.Stdout)
    }

    if strings.Count(resp.Stdout, "\n") != 1 {
        t.Fatalf("expected exactly one newline, got %d in:\n%q", strings.Count(resp.Stdout, "\n"), resp.Stdout)
    }
}
```
