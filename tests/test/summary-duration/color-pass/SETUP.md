# Scenario

**Feature**: --color grays inline duration and greens only the PASS token on final line

```
# color forced on
doctest test --color <dir> -> gray DURATION after inline parens -> green PASS(p/t) + plain " in DURATION"
```

## Preconditions

- A temp tree with 1 passing leaf exists.

## Steps

1. Create a 1-pass temp tree.
2. Run `doctest test --color <dir>`.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	testDir := createPassFailTree(t, 1, 0)
	req.Args = []string{"test", "--color", testDir}
	return nil
}
```