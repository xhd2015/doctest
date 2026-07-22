---
label: heavy
explanation: CLI edit contract via doctest binary
---

## Expected

- Exit 0; ASSERT.md frontmatter lists both labels comma-separated.

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
	got := readAssertFile(t, req.Args[1])
	wantPrefix := strings.Join([]string{
		"---",
		"label: ui-automation, manual",
		"explanation: first label",
		"---",
		"",
		"## Expected",
	}, "\n")
	if !strings.HasPrefix(got, wantPrefix) {
		t.Fatalf("ASSERT.md prefix mismatch\nwant prefix:\n%s\ngot:\n%s", wantPrefix, got)
	}
}
```