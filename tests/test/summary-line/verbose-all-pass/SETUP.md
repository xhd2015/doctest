# Scenario

**Feature**: verbose mode prints PASS summary after go test -v output

```
# verbose path
doctest test -v <dir> -> go test -v output -> PASS(passed/total)
```

## Preconditions

- A temp tree with 1 passing leaf exists.

## Steps

1. Create a 1-pass temp tree.
2. Run `doctest test -v <dir>`.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	testDir := createPassFailTree(t, 1, 0)
	req.Args = []string{"test", "-v", testDir}
	return nil
}
```