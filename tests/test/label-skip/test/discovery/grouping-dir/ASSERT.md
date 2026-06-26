## Expected

- Exit 0, PASS(1/1), exact skip block for labeled child under grouping dir.

```go
import (
	"path/filepath"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("exit code = %d\nstdout:\n%s\nstderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	treeRoot := filepath.Dir(req.Args[1])
	assertSkipBlockExact(t, resp.Stdout, treeRoot, "e2e/labeled_child", "ui-automation", "grouping skip")
	assertResultSummary(t, resp.Stdout, 1, 1)
}
```