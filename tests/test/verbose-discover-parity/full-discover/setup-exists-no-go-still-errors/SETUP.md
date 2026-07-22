# Scenario

**Feature**: when intermediate SETUP.md exists without a Go block, full discover still errors

```
createIntermediateSetupNoGoFixture  # intermediate/SETUP.md prose only
core.DiscoverTreeCases(parent)
  -> error containing intermediate + SETUP.md + must have a Go code block
```

## Preconditions

- SETUP.md **exists** on intermediate — still a real grouping contract violation.
- Expect **GREEN** before and after Fix 2 (behavior must not regress).

## Steps

1. Create fixture with prose-only intermediate SETUP.
2. Run full discover on parent.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	_, parentDir := createIntermediateSetupNoGoFixture(t)
	req.Op = "discover_full"
	req.DiscoverRoot = parentDir
	return nil
}
```
