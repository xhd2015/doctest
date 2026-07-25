# Scenario

**Feature**: removing a source leaf makes its gen package an orphan (`# deleted`)

```
run1: test tree  (leaves keep + gone)
rm -rf work/tree/gone
run2: test tree + gen-plan
  -> tree/gone/leaf.go gone
  -> # deleted
  -> not in manifest
```

Single-tree so tree-scope prune runs.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/libdoc/testtree"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	prepareSingleTreeModuleWithLeaves(t, req, []testtree.LeafSpec{
		{Name: "keep", Steps: "keep", Expected: "ok"},
		{Name: "gone", Steps: "gone", Expected: "ok"},
	})
	args := baseArgs(req, "tree")
	req.ArgsFull = args
	req.ArgsSubset = append([]string(nil), args...)
	req.RemoveSourceRels = []string{"tree/gone"}
	req.DebugEnv = debugGenPlanBypass
	return nil
}
```
