# Scenario

**Feature**: `-a` wipes gen root; second run recreates as mostly `# new`

```
run1: test tree (warm ledger)
run2: test -a tree + gen-plan
  -> new>=1 (cold-like rewrite after wipe)
```

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	prepareSingleTreeModule(t, req)
	req.ArgsFull = baseArgs(req, "tree")
	req.ArgsSubset = baseArgsForceA(req, "tree")
	req.DebugEnv = debugGenPlanBypass
	return nil
}
```
