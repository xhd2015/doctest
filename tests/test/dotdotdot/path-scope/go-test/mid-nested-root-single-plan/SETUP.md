# Scenario

**Case (S1):** `./tree/mid/...` with mid leaf + nested DOCTEST under mid (same module).

**Expected plan after fix:** exactly one go test under the gen root:

```text
cd <gen> && go test -v -count=1 ./tree/mid/...
```

Parent `...` already covers `tree/mid/nested/...` — no second `./tree/mid/nested/...` cmd.
Nested and mid leaves each run **once**.

```go
import (
	"os"
	"path/filepath"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.WorkDir = createPathScopeMidNestedRoot(t)
	genDir := filepath.Join(req.WorkDir, ".gen")
	if err := os.MkdirAll(genDir, 0755); err != nil {
		t.Fatal(err)
	}
	req.Args = []string{"test", "-v", "--label-all", "--gen-dir", genDir, "-count=1", "./tree/mid/..."}
	return nil
}
```
