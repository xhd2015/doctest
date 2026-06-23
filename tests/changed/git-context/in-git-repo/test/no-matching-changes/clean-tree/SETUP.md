# Scenario

**Feature**: clean working tree yields no-tests-changed warning

```
# no modifications
clean tree -> doctest test --changed -> warning, exit 0
```

## Steps

1. Create flat two-leaf tree and commit with no further changes.
2. Run `doctest test --changed`.

```go
import (
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	fx := createFlatTwoLeafTree(t)
	req.WorkDir = fx.RepoDir
	req.Args = []string{"test", fx.TreeDir, "--changed"}
	return nil
}
```