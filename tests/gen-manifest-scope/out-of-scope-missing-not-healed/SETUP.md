# Scenario

**Feature**: out-of-scope managed missing file is not healed by subset; healed when in scope

```
run1: test tree-a tree-b
rm G/tree-b/b1/leaf.go
run2: test tree-a only
  -> tree-b leaf still missing; still in manifest; not # new in run2
run3: test tree-b
  -> leaf recreated as # new
```

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	prepareTwoTreeModule(t, req)
	req.ArgsFull = baseArgs(req, "tree-a", "tree-b")
	req.ArgsSubset = baseArgs(req, "tree-a")
	req.ArgsThird = baseArgs(req, "tree-b")
	req.DeleteGenRels = []string{"tree-b/b1/leaf.go"}
	req.DebugEnv = debugGenPlanBypass
	return nil
}
```
