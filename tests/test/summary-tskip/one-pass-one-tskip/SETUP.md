# Scenario

**Feature**: one pass + one runtime t.Skip yields PASS (1/1, 1 t.Skip) and exit 0

```
# quiet path core contract
doctest test --no-color <ok + skip_me>
  -> suite-leaf Actions: 1 pass, 1 skip
  -> final summary PASS (1/1, 1 t.Skip) in …
  -> exit 0 (skips alone do not fail)
```

## Preconditions

- Root Setup has set `req.Bin` and a generous outer timeout.
- Fixture has two leaves: always-pass `ok` and `skip_me` that calls `t.Skip`.

## Steps

1. Create the one-pass-one-tskip fixture tree.
2. Isolate caches for the harness invocation.
3. Run `doctest test --no-color <fixture>` (quiet).

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	tree := createOnePassOneTSkipTree(t)
	req.Args = []string{"test", "--no-color", tree}
	req.Env = append(req.Env, isolateRunEnv(t)...)
	return nil
}
```
