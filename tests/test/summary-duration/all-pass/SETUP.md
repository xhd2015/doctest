# Scenario

**Feature**: three passing leaves report duration in inline and final summaries

```
# discover leaves, run packages, measure wall time
doctest test <dir> -> discover 3 leaves -> inline (3 Run, ...) in DURATION -> PASS(3/3) in DURATION
```

## Preconditions

- A temp tree with 3 passing leaves exists.

## Steps

1. Create a 3-pass / 0-fail temp tree.
2. Run `doctest test <dir>` (non-verbose).

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	testDir := createPassFailTree(t, 3, 0)
	req.Args = []string{"test", testDir}
	return nil
}
```