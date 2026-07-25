# Scenario

**Feature**: deleting gen go.mod is repaired on next generate

```
run1: test tree
rm G/go.mod
run2: same
  -> G/go.mod exists again
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
	req.DeleteGenRels = []string{"go.mod"}
	req.DebugEnv = debugGenPlanBypass
	return nil
}
```
