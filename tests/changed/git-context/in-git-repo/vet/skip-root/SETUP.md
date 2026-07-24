# Scenario

**Feature**: vet omits unchanged root `DOCTEST.md` even when invalid

```
# root missing ## Version (invalid) but not in changed list
# only leaf_a ASSERT.md changed
ChangedDoctestMarkdownFiles -> [leaf_a/ASSERT.md] (no DOCTEST.md)
```

## Steps

1. Create tree with invalid root `DOCTEST.md` (no Version).
2. Set changed path to `leaf_a/ASSERT.md` only.
3. Assert root is not in the markdown list.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	fx := createVetSkipRootTree(t)
	applyPolicyBase(req, fx)
	req.Policy = PolicyVetMD
	req.ChangedFiles = []string{treeRel(fx, "leaf_a", "ASSERT.md")}
	return nil
}
```
