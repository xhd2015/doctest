# Scenario

**Feature**: `doctest vet --changed` errors without git

```
# no git repo
doctest vet --changed <dir> -> error, non-zero exit
```

## Steps

1. Run `doctest vet --changed` on the non-git fixture tree.

```go
import (
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	treeDir := filepath.Join(req.WorkDir, "tests")
	req.Args = []string{"vet", treeDir, "--changed"}
	return nil
}
```