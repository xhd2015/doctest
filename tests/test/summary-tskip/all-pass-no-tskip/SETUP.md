# Scenario

**Feature**: all-pass with zero runtime skips keeps PASS (N/N) without t.Skip text

```
# regression: no runtime skip
doctest test --no-color <2 always-pass>
  -> suite-leaf Actions: 2 pass, 0 skip
  -> final summary PASS (2/2) in …  (no ", N t.Skip")
  -> exit 0
```

## Preconditions

- Root Setup has set `req.Bin` and a generous outer timeout.
- Fixture is two always-pass leaves (no `t.Skip`).

## Steps

1. Create a 2-pass / 0-fail temp tree.
2. Isolate caches for the harness invocation.
3. Run `doctest test --no-color <fixture>` (quiet).

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	tree := createAllPassNoTSkipTree(t)
	req.Args = []string{"test", "--no-color", tree}
	req.Env = append(req.Env, isolateRunEnv(t)...)
	return nil
}
```
