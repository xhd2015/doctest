# Scenario

**Feature**: single path arg plan uses arg[1/1] and includes bookkeeping in that tree

```
doctest test <tree> --gen-dir G
  -> gen-plan: arg[1/1] <tree>
       go.mod / go.sum / doctest.gen-manifest / doctest.tidy-done
       <package hierarchy>
  # no separate gen-plan: merged
```

## Preconditions

- One DOCTEST root fixture.

## Steps

1. prepareSingleTree; Args = test --gen-dir --no-color tree.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	prepareSingleTree(t, req)
	req.Args = baseTestArgs(req, "--no-color", req.TreeRoot)
	return nil
}
```
