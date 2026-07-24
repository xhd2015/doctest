# Scenario

**Feature**: quiet `doctest test` on parent with nested-only intermediate succeeds (1 case)

```
# quiet path already correct
createMegaParentNestedFixture
doctest test parent/ -> exit 0; planned 1 parent leaf
```

## Preconditions

- Intermediate has no SETUP; only nested DOCTEST child under intermediate/.
- Parent has one own leaf (`own_leaf`).

## Steps

1. Create mega parent nested fixture.
2. Run `doctest test --no-color <parentDir>`.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_, parentDir := createMegaParentNestedFixture(t)
	req.Args = []string{"test", "--no-color", parentDir}
	return nil
}
```
