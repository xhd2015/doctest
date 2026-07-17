---
label: heavy
---

## Expected
- Header/summary block (lines with `═══`, `Session:`, `Status:`, etc.) has NO timestamp prefix.
- Event display lines (with `  [%d] ` pattern) HAVE a timestamp prefix.
- Event continuation lines (`       ...`) HAVE a timestamp prefix.
- stderr is empty.

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
    if resp.ExitCode != 0 {
        t.Fatalf("exit code = %d, want 0\nstderr:\n%s", resp.ExitCode, resp.Stderr)
    }

    tsPrefix := regexp.MustCompile(`^\[\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\]`)
    eventPattern := regexp.MustCompile(`^\[\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\]   \[\d+\] `)
    continuationPattern := regexp.MustCompile(`^\[\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\]        `)

    lines := strings.Split(resp.Stdout, "\n")
    eventCount := 0

    for _, line := range lines {
        trimmed := strings.TrimSpace(line)
        if trimmed == "" {
            continue
        }

        hasTS := tsPrefix.MatchString(trimmed)
        isEvent := eventPattern.MatchString(trimmed)
        isContinuation := continuationPattern.MatchString(trimmed)
        isHeader := strings.Contains(trimmed, "═══") ||
            strings.HasPrefix(trimmed, "Session:") ||
            strings.HasPrefix(trimmed, "Status:") ||
            strings.HasPrefix(trimmed, "Runner:") ||
            strings.HasPrefix(trimmed, "Created:") ||
            strings.HasPrefix(trimmed, "Codex:") ||
            strings.HasPrefix(trimmed, "Opencode:") ||
            strings.HasPrefix(trimmed, "Events:")

        if isEvent || isContinuation {
            eventCount++
            if !hasTS {
                t.Fatalf("event/continuation line missing timestamp prefix:\n%q", line)
            }
        } else if isHeader {
            if hasTS {
                t.Fatalf("header line has unexpected timestamp prefix:\n%q", line)
            }
        } else if hasTS && !strings.Contains(trimmed, "No events yet") {
            t.Fatalf("unexpected timestamped line:\n%q", line)
        }
    }

    if eventCount == 0 {
        t.Fatal("no event display lines found in output")
    }
}
```
