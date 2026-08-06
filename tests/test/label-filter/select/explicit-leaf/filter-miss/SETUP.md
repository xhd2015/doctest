# Scenario

**Feature**: explicit leaf that fails label filter is skipped

```
SubDir=slow + --label e2e → skip {slow} reason label filter
```

## Steps

1. Fixture mod; SubDir = slow; LabelExprs = ["flaky"].

```go
import (
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.TreeRoot = writeLabelFilterMod(t)
	req.SubDir = filepath.Join(req.TreeRoot, "slow")
	req.LabelExprs = []string{"flaky"}
	req.ExplicitLeaf = true
	return nil
}
```
