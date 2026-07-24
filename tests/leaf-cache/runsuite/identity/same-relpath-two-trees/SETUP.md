# Scenario

**Feature**: same relative leaf path under two trees yields distinct identities

```
treeA + treeB, leafRel="leaf"
FormatLeafIdentity(A,"leaf") != FormatLeafIdentity(B,"leaf")
```

## Preconditions

- Twin trees from parent; LeafDir / LeafDirB set.

## Steps

1. Op=`format_identity` (already set).
2. Assert Identity and Identity2 differ and are non-empty.

## Context

- Bare relative paths would collide in multi-tree skip/fail maps.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	req.Op = "format_identity"
	// Twin trees already prepared by parent.
	return nil
}
```
