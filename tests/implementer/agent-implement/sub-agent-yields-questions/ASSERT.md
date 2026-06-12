## Expected
- Exit code 0 (yielding questions is not an error).
- Stdout contains the agent's response, the QUESTIONS separator, and the question JSON with options.

## Exit Code
- Exit code 0.

```go
import (
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
    if !strings.Contains(resp.Stdout, "working on it") {
        t.Fatalf("stdout missing agent response:\n%s", resp.Stdout)
    }
    if !strings.Contains(resp.Stdout, "QUESTIONS") {
        t.Fatalf("stdout missing QUESTIONS separator:\n%s", resp.Stdout)
    }
    if !strings.Contains(resp.Stdout, `"question"`) {
        t.Fatalf("stdout missing question:\n%s", resp.Stdout)
    }
    if !strings.Contains(resp.Stdout, `What is the target port?`) {
        t.Fatalf("stdout missing question text:\n%s", resp.Stdout)
    }
    if !strings.Contains(resp.Stdout, `"option"`) {
        t.Fatalf("stdout missing option:\n%s", resp.Stdout)
    }
    if !strings.Contains(resp.Stdout, `"explanation"`) {
        t.Fatalf("stdout missing explanation:\n%s", resp.Stdout)
    }
    if !strings.Contains(resp.Stdout, `default development port`) {
        t.Fatalf("stdout missing option explanation text:\n%s", resp.Stdout)
    }
}
```
