# Scenario

**Feature**: warm second generate with same args has no managed deletes

```
run1: test tree --gen-dir G  (cold)
run2: test tree --gen-dir G  (warm, identical)
  -> gen-plan summary deleted=0
  -> tree/* still in manifest
```

## Preconditions

- Single-tree fixture so tree-scope prune actually runs on warm path.
- `DOCTEST_DEBUG=gen-plan=1,bypass-go-test=1` for summary.

## Steps

1. prepareSingleTreeModule.
2. ArgsFull == ArgsSubset == `tree`.
3. DebugEnv enables gen-plan + bypass.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	prepareSingleTreeModule(t, req)
	args := baseArgs(req, "tree")
	req.ArgsFull = args
	req.ArgsSubset = append([]string(nil), args...)
	req.DebugEnv = "gen-plan=1,bypass-go-test=1"
	return nil
}
```
