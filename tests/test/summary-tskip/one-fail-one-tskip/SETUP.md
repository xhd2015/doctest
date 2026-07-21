# Scenario

**Feature**: one fail + one runtime t.Skip yields FAIL (0/1, 1 t.Skip) and non-zero exit

```
# quiet path fail + skip
doctest test --no-color <z_fail + skip_me>
  -> suite-leaf Actions: 1 fail, 1 skip
  -> final summary FAIL (0/1, 1 t.Skip) in …
  -> exit ≠ 0
```

## Preconditions

- Root Setup has set `req.Bin` and a generous outer timeout.
- Fixture has two leaves: forced-fail `z_fail` and `skip_me` that calls `t.Skip`.

## Steps

1. Create the one-fail-one-tskip fixture tree.
2. Isolate caches for the harness invocation.
3. Run `doctest test --no-color <fixture>` (quiet).

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	tree := createOneFailOneTSkipTree(t)
	req.Args = []string{"test", "--no-color", tree}
	req.Env = append(req.Env, isolateRunEnv(t)...)
	return nil
}
```
