## Expected
- Exit code 0.
- Header/border lines (with `═══`) have NO timestamp prefix `[...]`.
- `(no events yet)` status message HAS a timestamp prefix.
- `Done (session finished)` message HAS a timestamp prefix.
- Separator lines (`───`) have NO timestamp prefix.
- Stderr is empty.

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

    lines := strings.Split(resp.Stdout, "\n")
    foundNoEvents := false
    foundDone := false

    for _, line := range lines {
        trimmed := strings.TrimSpace(line)
        if trimmed == "" {
            continue
        }

        hasTS := tsPrefix.MatchString(trimmed)
        hasBorder := strings.Contains(trimmed, "═══") || strings.Contains(trimmed, "───")

        if strings.Contains(trimmed, "(no events yet)") {
            foundNoEvents = true
            if !hasTS {
                t.Fatalf("'(no events yet)' line missing timestamp prefix:\n%q", line)
            }
        } else if strings.Contains(trimmed, "Done (session finished)") {
            foundDone = true
            if !hasTS {
                t.Fatalf("'Done (session finished)' line missing timestamp prefix:\n%q", line)
            }
        } else if hasBorder {
            if hasTS {
                t.Fatalf("border line has unexpected timestamp prefix:\n%q", line)
            }
        } else if hasTS && !strings.Contains(trimmed, "no events") && !strings.Contains(trimmed, "Done") {
            t.Fatalf("unexpected timestamped line:\n%q", line)
        }
    }

    if !foundNoEvents {
        t.Fatal("missing '(no events yet)' message")
    }
    if !foundDone {
        t.Fatal("missing 'Done (session finished)' message")
    }
}
```
