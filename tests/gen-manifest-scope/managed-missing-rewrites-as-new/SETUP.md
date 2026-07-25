# Scenario

**Feature**: deleting a managed gen file then re-running rewrites it as `# new`

```
run1: test tree --gen-dir G
rm G/tree/leaf/leaf.go
run2: same + gen-plan
  -> leaf.go back on disk
  -> # new (or summary new>=1)
  -> still in manifest
```

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	prepareSingleTreeModule(t, req)
	args := baseArgs(req, "tree")
	req.ArgsFull = args
	req.ArgsSubset = append([]string(nil), args...)
	req.DeleteGenRels = []string{"tree/leaf/leaf.go"}
	req.DebugEnv = debugGenPlanBypass
	return nil
}
```
