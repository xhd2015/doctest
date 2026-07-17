## Expected

- Exit code 0.
- stdout contains `group/slow-leaf` before `group/mid-leaf` (index order).
- stdout contains `group/fast-leaf` after the slower paths (or at least after slow-leaf).

## Exit Code

- 0

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("exit=%d stderr=%s stdout=%s", resp.ExitCode, resp.Stderr, resp.Stdout)
	}
	out := resp.Stdout
	slow := strings.Index(out, "group/slow-leaf")
	mid := strings.Index(out, "group/mid-leaf")
	fast := strings.Index(out, "group/fast-leaf")
	if slow < 0 {
		t.Fatalf("missing slow-leaf in top output:\n%s", out)
	}
	if mid >= 0 && mid < slow {
		t.Fatalf("mid-leaf appeared before slow-leaf:\n%s", out)
	}
	if fast >= 0 && fast < slow {
		t.Fatalf("fast-leaf appeared before slow-leaf:\n%s", out)
	}
}
```
