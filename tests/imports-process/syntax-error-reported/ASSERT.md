---
label: heavy
---

## Expected
- Exit code is non-zero (doctest test fails).
- The stderr contains an error message related to formatting (`imports.Process` failure).
- No corrupted output file is written (the original temp file was cleaned up).

```go
import (
    "strings"
    "testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
    if err != nil {
        t.Fatalf("run error: %v", err)
    }

    if resp.ExitCode == 0 {
        t.Fatalf("expected non-zero exit code for syntax error, got %d\nstdout: %s\nstderr: %s", resp.ExitCode, resp.Stdout, resp.Stderr)
    }

    combined := resp.Stdout + "\n" + resp.Stderr
    hasFormatError := strings.Contains(combined, "format") ||
        strings.Contains(combined, "imports") ||
        strings.Contains(combined, "syntax") ||
        strings.Contains(combined, "error")
    if !hasFormatError {
        t.Fatalf("expected an error message about formatting/syntax, got:\nstdout: %s\nstderr: %s", resp.Stdout, resp.Stderr)
    }
}
```
