---
explanation: L3 process smoke — test --changed requires git repository
---

## Expected

- Non-zero exit code.
- stderr indicates `--changed` requires a git repository.

## Errors

- Command must fail when not inside a git work tree.

## Exit Code

Non-zero

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}
	if resp.ExitCode == 0 {
		t.Fatalf("expected non-zero exit, got 0\nstdout:\n%s\nstderr:\n%s", resp.Stdout, resp.Stderr)
	}
	lower := strings.ToLower(resp.Stderr)
	if !strings.Contains(lower, "git") {
		t.Fatalf("stderr should mention git, got:\n%s", resp.Stderr)
	}
}
```
