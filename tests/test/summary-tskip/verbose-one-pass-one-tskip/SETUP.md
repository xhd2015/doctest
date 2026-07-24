# Scenario

**Feature**: verbose (-v) path shows the same PASS (1/1, 1 t.Skip) as quiet

```
# always-json / -v presentation parity
doctest test -v --no-color <ok + skip_me>
  -> same suite-leaf counts as quiet
  -> final summary PASS (1/1, 1 t.Skip) in …
  -> exit 0
```

## Preconditions

- Root Setup has set `req.Bin` and a generous outer timeout.
- Same fixture shape as `one-pass-one-tskip` (1 pass + 1 t.Skip).

## Steps

1. Create the one-pass-one-tskip fixture tree.
2. Isolate caches for the harness invocation.
3. Run `doctest test -v --no-color <fixture>`.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	tree := createOnePassOneTSkipTree(t)
	req.Args = []string{"test", "-v", "--no-color", tree}
	req.Env = append(req.Env, isolateRunEnv(t)...)
	return nil
}
```
