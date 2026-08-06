# Scenario

**Feature**: quiet and `-v` prepare the same parent-tree case set when intermediate only holds a nested DOCTEST

```
# mega parent fixture
parent/own_leaf/              -> 1 parent case
parent/intermediate/          -> no SETUP, no ASSERT
parent/intermediate/nested/   -> nested DOCTEST (firewall)

# quiet path (Light→Hydrate)
doctest test parent/ -> exit 0; planned 1

# verbose path must match
doctest test -v parent/ -> exit 0; planned 1
  -/-> intermediate/SETUP.md: must have a Go code block
```

## Preconditions

- Shared mega fixture writer from root SETUP.
- Nested CLI runs use `e2e` when full integration.

## Steps

1. Create mega parent+nested intermediate fixture per leaf.
2. Run quiet and/or verbose `doctest test` against the parent root.

## Context

- Case selection remains Light → label filter → Hydrate only.
- `-v` must not call full discover as a **hard fail** after Light succeeded.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	// Leaves set Args / dual_cli fields; ensure CLI op default.
	if req.Op == "" {
		req.Op = "cli"
	}
	return nil
}
```
