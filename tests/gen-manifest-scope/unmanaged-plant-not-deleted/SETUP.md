# Scenario

**Feature**: hand-planted file under gen package is not treated as orphan

```
run1: test tree --gen-dir G
plant: G/tree/__droot/unused.go   # not via WriteIfChanged
run2: test tree --gen-dir G
  -> unused.go still on disk
  -> unused.go not in doctest.gen-manifest
  -> gen-plan deleted=0 (not listed as # deleted)
```

## Preconditions

- Single-tree fixture (tree prune runs; still must not delete untracked file).
- Plant path under a package that exists after generate (`tree/__droot/`).

## Steps

1. prepareSingleTreeModule.
2. Same args both runs; PlantRel = `tree/__droot/unused.go`.
3. gen-plan + bypass for deleted count observability.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	prepareSingleTreeModule(t, req)
	args := baseArgs(req, "tree")
	req.ArgsFull = args
	req.ArgsSubset = append([]string(nil), args...)
	req.PlantRel = "tree/__droot/unused.go"
	req.DebugEnv = "gen-plan=1,bypass-go-test=1"
	return nil
}
```
