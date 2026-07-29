# Scenario

**Feature**: single-gen run uses **one** workspace suite go test (ModeWorkspaceSuite)

Production path (Phase 1): single-cmd only — no multi-cmd path-shaped merge.

```
doctest test -v ./...  →  exactly one plan:
  cd <gen> && go test ./__workspace/suite
```

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.WorkDir = createSingleModTwoTrees(t)
	req.Args = []string{"test", "-v", "./..."}
	return nil
}
```
