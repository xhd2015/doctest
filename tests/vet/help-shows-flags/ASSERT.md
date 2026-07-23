---
explanation: L2 in-process CLI for vet --help argv/usage wiring
---

## Expected

- Exit code 0.
- stdout contains `-v`, `--verbose`.
- stdout contains `<dir...>` (multiple positional args).
- stdout contains `./...` examples.

## Exit Code

- 0

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("run failed unexpectedly: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected zero exit, got %d\nstdout:\n%s\nstderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	if !strings.Contains(resp.Stdout, "-v") {
		t.Fatalf("stdout missing -v flag:\n%s", resp.Stdout)
	}
	if !strings.Contains(resp.Stdout, "--verbose") {
		t.Fatalf("stdout missing --verbose flag:\n%s", resp.Stdout)
	}
	if !strings.Contains(resp.Stdout, "<dir...>") {
		t.Fatalf("stdout missing <dir...> positional pattern:\n%s", resp.Stdout)
	}
	if !strings.Contains(resp.Stdout, "./...") {
		t.Fatalf("stdout missing ./... example:\n%s", resp.Stdout)
	}
}
```
