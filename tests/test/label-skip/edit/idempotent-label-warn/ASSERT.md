## Expected

- Exit 0, exact stderr warning, ASSERT.md frontmatter unchanged.

```go
import (
	"path/filepath"
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
	leafDir := req.Args[1]
	wantWarn := wantIdempotentLabelWarning(leafDir, "ui-automation")
	assertStderrExact(t, resp.Stderr, wantWarn)
	got := readAssertFile(t, leafDir)
	wantPrefix := strings.Join([]string{
		"---",
		"label: ui-automation",
		"explanation: heavy ui test",
		"---",
		"",
		"## Expected",
	}, "\n")
	if !strings.HasPrefix(got, wantPrefix) {
		t.Fatalf("ASSERT.md must be unchanged\nwant prefix:\n%s\ngot:\n%s", wantPrefix, got)
	}
	_ = filepath.Base(leafDir)
}
```