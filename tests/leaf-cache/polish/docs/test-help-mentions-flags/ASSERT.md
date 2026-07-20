---
label: heavy
explanation: builds selftest binary to print test --help
---

## Expected

- Exit 0.
- Help text (stdout or stderr) contains `-a` and `--no-leaf-cache`.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("help exit %d\nstdout:\n%s\nstderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	out := resp.Stdout + "\n" + resp.Stderr
	hasA := strings.Contains(out, "-a")
	hasNoLeaf := strings.Contains(out, "--no-leaf-cache")
	if !hasA || !hasNoLeaf {
		t.Fatalf("help missing leaf-cache flags (-a=%v no-leaf-cache=%v)\n%s", hasA, hasNoLeaf, out)
	}
}
```
