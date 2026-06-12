## Expected
- Exit code 0.
- The file contains a JSON line with `type` = `"question"` and the correct text.

```go
import (
    "os"
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

    qFile := os.Getenv("TEST_Q_FILE")
    if qFile == "" {
        t.Fatal("TEST_Q_FILE not set")
    }

    data, readErr := os.ReadFile(qFile)
    if readErr != nil {
        t.Fatalf("read questions file: %v", readErr)
    }
    content := string(data)
    if content == "" {
        t.Fatal("questions file is empty")
    }
    if !strings.Contains(content, `"question"`) {
        t.Fatalf("missing type field:\n%s", content)
    }
    if !strings.Contains(content, "What is the target port?") {
        t.Fatalf("missing question text:\n%s", content)
    }
}
```
