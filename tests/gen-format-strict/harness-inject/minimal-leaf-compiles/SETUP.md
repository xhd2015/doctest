# Scenario

**Feature**: A4 — minimal leaf relying on harness inject compiles when author only imports testing

```
leaf SETUP: import "testing"; Setup uses t and req only
  -> RunTestLeaf + session.Doctest inject
  -> suite go test succeeds
```

## Preconditions

- Should stay GREEN across Phase A (regression guard for harness inject).

## Steps

1. Inherit FixtureKind `harness-minimal`.
2. Assert RunErr empty.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Op = "run_fixture"
	if req.FixtureKind == "" {
		req.FixtureKind = "harness-minimal"
	}
	return nil
}
```
