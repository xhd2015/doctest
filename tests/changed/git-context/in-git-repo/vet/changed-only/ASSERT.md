---
label: heavy
---

## Expected

- Non-zero exit code.
- stderr contains embedded-Go anti-pattern error for the changed leaf.
- stderr does not reference the unchanged sibling `leaf_a` SETUP.

## Errors

- Vet must reject the anti-pattern in the changed `SETUP.md`.

## Exit Code

Non-zero

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}
	if resp.ExitCode == 0 {
		t.Fatalf("expected non-zero exit, got 0\nstdout:\n%s\nstderr:\n%s", resp.Stdout, resp.Stderr)
	}
	if !strings.Contains(resp.Stderr, "anti-pattern: raw Go code embedded in string literal") {
		t.Fatalf("stderr missing anti-pattern error:\n%s", resp.Stderr)
	}
	if strings.Contains(resp.Stderr, "leaf_a") {
		t.Fatalf("stderr should not mention unchanged sibling leaf_a:\n%s", resp.Stderr)
	}
}
```