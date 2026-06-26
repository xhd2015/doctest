## Expected

- Exit 0, PASS(1/1), exact skip block for labeled leaf discovered via `...`.

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
	wantBlock := "SKIPPED 1 TESTS\n  mod/labeled_leaf\n    label: ui-automation\n    explanation: dotdotdot skip"
	gotBlock := skipBlock(resp.Stdout)
	if gotBlock != wantBlock {
		t.Fatalf("skip block mismatch\nwant:\n%s\ngot:\n%s\nstdout:\n%s", wantBlock, gotBlock, resp.Stdout)
	}
	assertResultSummary(t, resp.Stdout, 1, 1)
}
```