# Scenario

**Feature**: full discover succeeds on parent whose intermediate dir only holds a nested DOCTEST (no SETUP)

```
createMegaParentNestedFixture
core.DiscoverTreeCases(parent)
  -> err == nil
  -> len(cases) == 1  # own_leaf only; nested firewall skipped
```

## Preconditions

- Intermediate has **no** SETUP.md file.
- Classic TDD: **RED** until full discover skips missing intermediate SETUP.

## Steps

1. Create mega parent nested fixture.
2. Set `DiscoverRoot` to parent path; `Op=discover_full`.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_, parentDir := createMegaParentNestedFixture(t)
	req.Op = "discover_full"
	req.DiscoverRoot = parentDir
	return nil
}
```
