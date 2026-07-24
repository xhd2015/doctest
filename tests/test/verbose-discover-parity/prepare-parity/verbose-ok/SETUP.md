# Scenario

**Feature**: `doctest test -v` on parent with nested-only intermediate must prepare OK

```
# same fixture as quiet-ok
createMegaParentNestedFixture
doctest test -v parent/
  -> exit 0
  -/-> intermediate/SETUP.md: must have a Go code block
  -> planned 1 (same as quiet)
```

## Preconditions

- Classic TDD: **RED** until `-v` no longer hard-fails on full discover re-walk
  (and/or full discover accepts pure nested intermediate without SETUP).

## Steps

1. Create mega parent nested fixture.
2. Run `doctest test -v --no-color <parentDir>`.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_, parentDir := createMegaParentNestedFixture(t)
	req.Args = []string{"test", "-v", "--no-color", parentDir}
	return nil
}
```
