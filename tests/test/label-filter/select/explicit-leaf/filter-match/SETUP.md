# Scenario

**Feature**: explicit labeled leaf runs when expression matches

```
SubDir=slow + --label slow → run {slow}
```

## Steps

1. Fixture mod; SubDir = slow leaf; LabelExprs = ["slow"].

```go
import (
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.TreeRoot = writeLabelFilterMod(t)
	req.SubDir = filepath.Join(req.TreeRoot, "slow")
	req.LabelExprs = []string{"slow"}
	req.ExplicitLeaf = true
	return nil
}
```
