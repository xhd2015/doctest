## Expected
- Header/summary block (lines with `═══`, `Session:`, `Status:`, `Runner:`, `Created:`, `Events:`) has NO timestamp prefix.
- `No events yet` message HAS a timestamp prefix.
- stderr is empty.

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
    if resp.ExitCode != 0 {
        t.Fatalf("exit code = %d, want 0\nstderr:\n%s", resp.ExitCode, resp.Stderr)
    }

    tsPrefix := regexp.MustCompile(`^\[\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\]`)

    lines := strings.Split(resp.Stdout, "\n")
    foundNoEvents := false

    for _, line := range lines {
        trimmed := strings.TrimSpace(line)
        if trimmed == "" {
            continue
        }

        hasTS := tsPrefix.MatchString(trimmed)
        isHeader := strings.Contains(trimmed, "═══") ||
            strings.HasPrefix(trimmed, "Session:") ||
            strings.HasPrefix(trimmed, "Status:") ||
            strings.HasPrefix(trimmed, "Runner:") ||
            strings.HasPrefix(trimmed, "Created:") ||
            strings.HasPrefix(trimmed, "Codex:") ||
            strings.HasPrefix(trimmed, "Opencode:") ||
            strings.HasPrefix(trimmed, "Events:")

        if strings.Contains(trimmed, "No events yet") {
            foundNoEvents = true
            if !hasTS {
                t.Fatalf("'No events yet' line missing timestamp prefix:\n%q", line)
            }
        } else if isHeader {
            if hasTS {
                t.Fatalf("header line has unexpected timestamp prefix:\n%q", line)
            }
        } else if hasTS {
            t.Fatalf("unexpected timestamped line:\n%q", line)
        }
    }

    if !foundNoEvents {
        t.Fatal("missing 'No events yet' message")
    }
}
```
