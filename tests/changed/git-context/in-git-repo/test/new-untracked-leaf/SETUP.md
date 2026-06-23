# Scenario

**Feature**: an untracked new leaf is the only one that runs

```
# new leaf directory added but not committed
untracked leaf_c/ -> doctest test --changed -> 1 Run (new leaf only)
```

## Steps

1. Create flat two-leaf tree and commit.
2. Add a new `leaf_c` directory (untracked).
3. Run `doctest test --changed`.

```go
import (
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	fx := createFlatTwoLeafTree(t)
	writeLeaf(t, fx.TreeDir, "leaf_c")
	req.WorkDir = fx.RepoDir
	req.Args = []string{"test", fx.TreeDir, "--changed"}
	return nil
}
```