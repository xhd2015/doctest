# Scenario

**Feature**: real outer leaf failures still produce FAIL (p/t) with non-zero exit

```
# one forced-fail outer leaf (no nested intentional pass-around)
doctest test --no-color <1-fail>
  -> go test fails the leaf
  -> final summary FAIL (0/1), exit ≠ 0
```

## Preconditions

- Root Setup has set `req.Bin` and a generous outer timeout.
- Fixture is a single Assert-fail leaf (regression guard for always-json).

## Steps

1. Create a temp 1-fail tree.
2. Isolate caches for the harness invocation.
3. Run `doctest test --no-color <tree>`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	testDir := createRealFailTree(t)
	req.Args = []string{"test", "--no-color", testDir}
	req.Env = append(req.Env, isolateRunEnv(t)...)
	return nil
}
```
