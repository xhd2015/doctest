# Scenario

**Feature**: FormatLeafIdentity is stable and distinguishes leaf rels under one tree

```
FormatLeafIdentity(tree, "leaf") twice -> equal
FormatLeafIdentity(tree, "leaf") != FormatLeafIdentity(tree, "other")
```

## Preconditions

- One tree root from twins (TreeRoot); synthetic second rel path `other`
  need not exist on disk — identity is a pure string function of roots + rel.

## Steps

1. Op=`format_identity_stable` — Run formats identity twice for `leaf` and once
   for a different rel under the same tree.
2. Assert Key==Key2 (stability) and Identity != Identity2 (rel distinction).

## Context

- Stability is required so prepare and record share map keys.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	t.Helper()
	req.Op = "format_identity_stable"
	// TreeRoot already set by parent prepareTwinTrees.
	return nil
}
```
