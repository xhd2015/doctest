# Scenario

**Feature**: explicit leaf that fails label filter is skipped

```
SubDir=slow + --label heavy → skip {slow} reason label filter
```

## Steps

1. Fixture mod; SubDir = slow; LabelExprs = ["heavy"].

```go
import (
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	req.TreeRoot = writeLabelFilterMod(t)
	req.SubDir = filepath.Join(req.TreeRoot, "slow")
	req.LabelExprs = []string{"heavy"}
	req.ExplicitLeaf = true
	return nil
}
```
