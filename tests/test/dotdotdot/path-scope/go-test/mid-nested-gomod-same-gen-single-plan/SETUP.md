# Scenario

**Case (S4):** `./tree/mid/...` with mid leaf + nested go.mod DOCTEST under mid.

Shared `--gen-dir` flattens both into one gen module. Rule:

> same `cd` → **one** `go test` with combined patterns  
> (or different `cd` if multi-gen — never same dir twice)

Expected (shared gen):

```text
cd <gen> && go test -v -count=1 ./suite/... ./tree/mid/...
```

```go
import (
	"os"
	"path/filepath"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.WorkDir = createPathScopeMidNestedGomod(t)
	genDir := filepath.Join(req.WorkDir, ".gen")
	if err := os.MkdirAll(genDir, 0755); err != nil {
		t.Fatal(err)
	}
	req.Args = []string{"test", "-v", "--label-all", "--gen-dir", genDir, "-count=1", "./tree/mid/..."}
	return nil
}
```
