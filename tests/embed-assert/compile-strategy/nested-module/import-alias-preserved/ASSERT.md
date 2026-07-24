## Expected

- `doctest test` exits 0.
- Generated `leaf_test.go` contains `outputassert` alias referencing assert module path.

## Exit Code

- Exit code 0.

```go
import (
	"os"
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("doctest subprocess failed: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected aliased assert import test to pass, got exit %d\nstdout:\n%s\nstderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	genTest := generatedLeafTestPath(req.OutsideGenDir)
	data, err := os.ReadFile(genTest)
	if err != nil {
		t.Fatalf("read generated test: %v", err)
	}
	src := string(data)
	if !strings.Contains(src, `outputassert "github.com/xhd2015/doctest/assert"`) {
		t.Fatalf("expected aliased assert import preserved, got:\n%s", src)
	}
	if !strings.Contains(src, "outputassert.Output") {
		t.Fatalf("expected outputassert.Output call preserved, got:\n%s", src)
	}
}
```