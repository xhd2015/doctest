# Scenario

**Feature**: full `DiscoverTreeCases` accepts pure nested-tree parent dirs; SETUP that exists still requires Go Setup

```
# pure nested intermediate (no SETUP file)
DiscoverTreeCases(parent) -> OK; cases = [own_leaf]

# intermediate SETUP.md exists without Go block
DiscoverTreeCases(parent) -> error: intermediate/SETUP.md: must have a Go code block
```

## Preconditions

- Leaves use `Op=discover_full` (in-process API; no nested CLI compile).
- Fixtures from root SETUP.

## Steps

1. Build the appropriate parent fixture.
2. Point `req.DiscoverRoot` at the parent tree root.

## Context

- Full discover must not invent a SETUP requirement when the file is absent.
- Real grouping still requires Go Setup when SETUP.md is present.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Op = "discover_full"
	return nil
}
```
