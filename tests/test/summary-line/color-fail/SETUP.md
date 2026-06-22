# Scenario

**Feature**: --color wraps entire FAIL summary in red

```
# color forced on with failure
doctest test --color <dir> -> red FAIL(passed/total)
```

## Preconditions

- A temp tree with 1 failing leaf exists.

## Steps

1. Create a 0-pass / 1-fail temp tree.
2. Run `doctest test --color <dir>`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	testDir := createPassFailTree(t, 0, 1)
	req.Args = []string{"test", "--color", testDir}
	return nil
}
```