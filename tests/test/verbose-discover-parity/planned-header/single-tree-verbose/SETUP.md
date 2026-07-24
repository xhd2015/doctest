# Scenario

**Feature**: `doctest test -v` on a single tree always prints planned test count before go test

```
createFastPassTree  # 1 leaf
doctest test -v --no-color <tree>
  -> user-visible planned: (1 tests) or ─── 1 test cases or 1 tests planned
  -> then go test -v
```

## Preconditions

- Classic TDD: **RED** until single-tree `-v` always surfaces planned N tests
  (not only a path line without count).

## Steps

1. Create 1-pass temp tree.
2. Run `doctest test -v --no-color <tree>`.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	testDir := createFastPassTree(t)
	req.Args = []string{"test", "-v", "--no-color", testDir}
	return nil
}
```
