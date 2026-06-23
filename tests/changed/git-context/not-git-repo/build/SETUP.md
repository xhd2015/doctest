# Scenario

**Feature**: `doctest build --changed` errors without git

```
# no git repo
doctest build --changed <dir> -> error, non-zero exit
```

## Steps

1. Run `doctest build --changed` on the non-git fixture tree.

```go
import (
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	treeDir := filepath.Join(req.WorkDir, "tests")
	genDir := filepath.Join(req.WorkDir, "gen")
	req.Args = []string{"build", treeDir, "--changed", "--gen-dir", genDir}
	return nil
}
```