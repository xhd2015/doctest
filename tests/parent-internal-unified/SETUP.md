# Scenario

**Feature**: parent-internal multi-leaf trees always use unified mapping-gen (P2)

```
# temp parent module with internal/greet + ≥2 leaves
createParentInternalMultiLeafModule
  -> runner.RunTest(tests, GenDir[, CoverProfile])
  -> unified suite (layout A), not multi-leaf internalCompile
```

## Preconditions

- Nested tree root is `tests/parent-internal-unified/` (`d.DOCTEST_ROOT`).
- No process `Setenv`/`Chdir`; temps via `t.TempDir()` and `d.DOCTEST_*`.
- P1 expose parent-internal eligibility already lands expose helpers; P2 flips
  generation mode to always unified for parent-internal cases.

## Steps

1. Root Setup is a thin default (Op defaults left to grouping/leaf).
2. `multi-leaf` Setup materializes the parent module + two subject leaves.
3. Leaf Setup may enable coverprofile; Assert checks pass / suite / layout / cover.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	_ = d
	// Defaults: no cover; multi-leaf Setup fills ModuleRoot/TestDir/GenDir.
	req.WithCover = false
	req.CoverPath = ""
	req.CoverPkg = ""
	req.CoverMode = ""
	req.ModuleRoot = ""
	req.TestDir = ""
	req.GenDir = ""
	return nil
}
```
