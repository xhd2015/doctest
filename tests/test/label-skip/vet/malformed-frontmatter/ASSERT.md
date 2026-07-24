---
label: heavy
explanation: CLI vet contract via doctest binary
---

## Expected

- Non-zero exit code.

```go
import (
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}
	if resp.ExitCode == 0 {
		t.Fatalf("expected non-zero exit for malformed frontmatter\nstdout:\n%s\nstderr:\n%s", resp.Stdout, resp.Stderr)
	}
}
```