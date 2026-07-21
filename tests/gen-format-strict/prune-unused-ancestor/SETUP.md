# Scenario

**Feature**: unused intermediate/ancestor package imports are pruned (or never emitted)

```
# hierarchy
root → feature (FeatureHelper) → mid (no FeatureHelper call) → leaf
  -> generate intermediate mid/setup.go
  -> unused parent (feature) import absent after WriteFormattedGo
```

## Preconditions

- Assemble may initially emit parent imports for intermediate packages.
- Write path must drop imports that no selector in the file uses.
- Leaves set FixtureKind `unused-parent`.

## Steps

1. `Op=run_fixture`, FixtureKind `unused-parent`.
2. After generate, read intermediate `mid/setup.go`.
3. Assert parent feature package import is absent.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Op = "run_fixture"
	req.FixtureKind = "unused-parent"
	return nil
}
```
