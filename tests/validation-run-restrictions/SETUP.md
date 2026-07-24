# Scenario

**Feature**: the doctest tree is discoverable

```
# enforce SETUP.md rules at build time
root missing Run -> build error
child redefines Run -> build error
non-root missing Setup -> build error
```

## Preconditions
- The doctest tree is discoverable.
- Testdata fixtures exist under `tests/testdata/` for invalid trees.

## Steps
1. Each leaf creates temporary doc-style test trees.
2. The root Run executes `doctest test` on the created tree.
3. Leaf Assert checks the outcome against the expected validation error or success.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    t.Log("validation-run-restrictions: testing Run placement rules")
    return nil
}
```
