# Scenario

**Feature**: under `-v`, user-visible output always includes planned trees/tests before go test

```
# single tree
doctest test -v <1-pass> -> planned line with N tests -> cd … && go test -v …

# workspace
doctest test -v ./... -> doctest: workspace (N trees, M tests) -> cd … && go test …
```

## Preconditions

- Nested CLI; use `e2e` when full integration.
- Quiet mode already prints planned counts (unchanged success path).

## Steps

1. Build single-tree or multi-tree fixtures.
2. Run `doctest test -v` and assert planned header appears.

## Context

- Workspace verbose currently prints only `cd … && go test` without the planned trees/tests line.
- Single-tree verbose omits `(N tests)` on the first announce line (quiet has it).

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	if req.Op == "" {
		req.Op = "cli"
	}
	return nil
}
```
