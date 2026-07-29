# Scenario

**Feature**: single-gen run uses workspace suite go test (today’s shape via gotestmap.Plan)

```
doctest test -v ./...  →  cd <gen> && go test ./__workspace/suite
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
