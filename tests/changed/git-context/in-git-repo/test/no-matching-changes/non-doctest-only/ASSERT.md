---
label: heavy
---

## Expected

- Exit code 0.
- stderr is empty (zero-change trees are silent without `-v`).

## Side Effects

- No tests are executed despite non-doctest file change.

## Exit Code

0

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
	if strings.TrimSpace(resp.Stderr) != "" {
		t.Fatalf("stderr should be empty without -v, got:\n%s", resp.Stderr)
	}
}
```