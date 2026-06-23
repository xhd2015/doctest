# Scenario

**Feature**: `doctest test --changed` errors without git

```
# no git repo
doctest test --changed <dir> -> error, non-zero exit
```

## Steps

1. Run `doctest test --changed` on the non-git fixture tree.

```go
import (
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	treeDir := filepath.Join(req.WorkDir, "tests")
	req.Args = []string{"test", treeDir, "--changed"}
	return nil
}
```