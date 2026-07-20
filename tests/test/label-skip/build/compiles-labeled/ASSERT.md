---
label: heavy
---

## Expected

- Exit 0; no test execution output.

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
		t.Fatalf("exit code = %d\nstdout:\n%s\nstderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	combined := resp.Stdout + resp.Stderr
	if strings.Contains(combined, "SKIPPED ") || strings.Contains(combined, "skipped ") {
		t.Fatalf("build must not print skip summary\ngot:\n%s", combined)
	}
	if strings.Contains(combined, "PASS (") {
		t.Fatalf("build must not print test PASS summary\ngot:\n%s", combined)
	}
}
```