---
label: heavy
---

## Expected

- Exit 0; ASSERT.md explanation field is exactly `first; second`.

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
	if !strings.Contains(got, "explanation: first; second") {
		t.Fatalf("expected appended explanation in ASSERT.md\ngot:\n%s", got)
	}
	wantFrontmatter := "---\nlabel: ui-automation\nexplanation: first; second\n---\n"
	if !strings.HasPrefix(got, wantFrontmatter) {
		t.Fatalf("frontmatter mismatch\nwant prefix:\n%s\ngot:\n%s", wantFrontmatter, got)
	}
}
```