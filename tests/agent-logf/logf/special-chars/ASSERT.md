---
label: heavy
---

## Expected
- Output has `[YYYY-MM-DDTHH:MM:SS] line1\nline2 -- special: !@#$%` followed by exactly one final newline.
- No double newline appended (message already ends with `\n`).
- Multiline content and special characters (`!`, `@`, `#`, `$`, `%`) are preserved.

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

    tsPattern := regexp.MustCompile(`^\[\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\] line1\nline2 -- special: !@#\$%\n$`)
    if !tsPattern.MatchString(resp.Stdout) {
        t.Fatalf("expected timestamped multiline message with special chars and exactly one trailing newline, got:\n%q", resp.Stdout)
    }

    if strings.Count(resp.Stdout, "\n") != 2 {
        t.Fatalf("expected exactly 2 newlines (1 internal + 1 trailing), got %d in:\n%q", strings.Count(resp.Stdout, "\n"), resp.Stdout)
    }
}
```
