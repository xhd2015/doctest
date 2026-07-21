# Scenario

**Feature**: A3 — intermediate setup.go does not keep unused parent package import

```
feature/SETUP exports FeatureHelper
mid/SETUP does not call FeatureHelper
  -> mid/setup.go after generate has no import of feature package
```

## Preconditions

- FixtureKind set by grouping parent (`unused-parent`).
- May still be RED if prune regresses when stdlib auto-add is removed.

## Steps

1. Inherit FixtureKind from grouping Setup.
2. Run generate (suite may pass).
3. Assert IntermediateSetupGo lacks unused parent import path.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	// FixtureKind already set by grouping; keep Op run_fixture.
	req.Op = "run_fixture"
	if req.FixtureKind == "" {
		req.FixtureKind = "unused-parent"
	}
	return nil
}
```
